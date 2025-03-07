package impl

import (
	"context"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/phongnd2802/ezy-mark/internal/cache"
	"github.com/phongnd2802/ezy-mark/internal/db"
	"github.com/phongnd2802/ezy-mark/internal/helpers"
	"github.com/phongnd2802/ezy-mark/internal/middlewares"
	"github.com/phongnd2802/ezy-mark/internal/models"
	"github.com/phongnd2802/ezy-mark/internal/pkg/crypto"
	"github.com/phongnd2802/ezy-mark/internal/pkg/random"
	"github.com/phongnd2802/ezy-mark/internal/pkg/token"
	"github.com/phongnd2802/ezy-mark/internal/pkg/utils"
	"github.com/phongnd2802/ezy-mark/internal/response"
	"github.com/phongnd2802/ezy-mark/internal/services"
	"github.com/phongnd2802/ezy-mark/internal/worker"
	"github.com/rs/zerolog/log"
)

const (
	OTP_EXPIRATION        int64 = 60
	ACCESS_TOKEN_TIME_EX  string = "1h"
	REFRESH_TOKEN_TIME_EX string = "168h"
)

type authServiceImpl struct {
	store       db.Store
	cache       cache.Cache
	distributor worker.TaskDistributor
}

func (s *authServiceImpl) Login(ctx context.Context, params *models.LoginRequest) (code int, res *models.LoginResponse, err error) {
	// Check Email
	userFound, err := s.store.GetUserBaseByEmail(ctx, params.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return response.ErrCodeAuthenticationFailed, nil, nil
		}
		return response.ErrCodeInternalServer, nil, err
	}

	matched := crypto.VerifyPassword(params.Password, userFound.UserPassword)
	if !matched {
		return response.ErrCodeAuthenticationFailed, nil, nil
	}

	// Check account active
	if !userFound.IsVerified.Bool {
		return response.ErrCodeAccountNotVerified, nil, nil
	}

	// Create UUID User
	subToken := utils.GenerateCliTokenUUID(userFound.UserID)
	log.Info().
		Str("subToken", subToken).
		Str("email", userFound.UserEmail).
		Int64("userID", userFound.UserID).
		Msg("User logged in")

	// Get UserProfile Profile
	profileUser, err := s.store.GetUserProfile(ctx, userFound.UserID)
	if err != nil {
		return response.ErrCodeInternalServer, nil, err
	}

	// Convert to Json
	profileUserJson, err := sonic.Marshal(profileUser)
	if err != nil {
		return response.ErrCodeInternalServer, nil, err
	}

	// Give profileUserJson to redsi with key = subToken
	duration, _ := time.ParseDuration(ACCESS_TOKEN_TIME_EX)

	err = s.cache.SetEx(ctx, subToken, profileUserJson, int64(duration.Seconds()))
	if err != nil {
		return response.ErrCodeInternalServer, nil, err
	}

	// Generate AccessToken and RefreshToken
	accessToken, err := token.CreateToken(subToken, ACCESS_TOKEN_TIME_EX)
	if err != nil {
		return response.ErrCodeInternalServer, nil, err
	}
	refreshToken, err := token.CreateToken(subToken, REFRESH_TOKEN_TIME_EX)
	if err != nil {
		return response.ErrCodeInternalServer, nil, err
	}

	// Update State Login
	userAgent, _ := ctx.Value(middlewares.UserAgentKey).(string)
	clientIP, _ := ctx.Value(middlewares.ClientIPKey).(string)
	go func() {
		newCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_, err := s.store.CreateUserSession(newCtx, db.CreateUserSessionParams{
			SessionID:    uuid.New(),
			RefreshToken: refreshToken,
			UserAgent:    userAgent,
			ClientIp:     clientIP,
			UserLoginTime: pgtype.Timestamptz{
				Time:             time.Now(),
				InfinityModifier: pgtype.Finite,
				Valid:            true,
			},
			UserID: userFound.UserID,
		})
		if err != nil {
			log.Error().AnErr("errUpdateLogin", err)
		}
	}()

	return response.ErrCodeSuccess, &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authServiceImpl) VerifyOTP(ctx context.Context, params *models.VerifyOTPReq) (int, *models.VerifyOTPRes, error) {
	userKeyOtp := helpers.GetUserKeyOtp(params.Token)
	userKeySession := helpers.GetUserKeySession(params.Token)
	exists, err := s.cache.Exists(ctx, userKeySession)
	if err != nil {
		return response.ErrCodeInternalServer, nil, err
	}
	if exists == 0 {
		return response.ErrCodeExpiredSession, nil, nil
	}
	otpFound, err := s.cache.Get(ctx, userKeyOtp)
	if err != nil {
		return response.ErrCodeInternalServer, nil, err
	}

	if otpFound != params.Otp {
		ttl, err := s.cache.TTL(ctx, userKeyOtp)
		if err != nil {
			return response.ErrCodeInternalServer, nil, err
		}
		return response.ErrCodeOtpDoesNotMatch, &models.VerifyOTPRes{
			TTL: int(ttl.Seconds()),
		}, nil
	}

	userUpdated, err := s.store.UpdateUserVerify(ctx, params.Token)
	if err != nil {
		return response.ErrCodeInternalServer, nil, err
	}

	err = s.cache.Del(ctx, userKeySession, userKeyOtp)
	if err != nil {
		return response.ErrCodeInternalServer, nil, err
	}

	nickName := helpers.ExtractNameFromEmail(userUpdated.UserEmail)
	_, err = s.store.CreateUserProfile(ctx, db.CreateUserProfileParams{
		UserID:       userUpdated.UserID,
		UserEmail:    userUpdated.UserEmail,
		UserNickname: nickName,
	})
	if err != nil {
		return response.ErrCodeInternalServer, nil, err
	}

	return response.ErrCodeSuccess, nil, nil
}

func (s *authServiceImpl) Register(ctx context.Context, params *models.RegisterRequest) (int, *models.RegisterResponse, error) {
	// Hash email
	hashEmail := crypto.GetHash(params.Email)
	log.Info().Str("HashEmail", hashEmail)
	// Check Email exists
	userFound, err := s.store.GetUserBaseByEmail(ctx, params.Email)
	if err != nil && err != pgx.ErrNoRows {
		return response.ErrCodeInternalServer, nil, err
	}

	if userFound.UserID > 0 {
		userKeySession := helpers.GetUserKeySession(hashEmail)
		ttl, err := s.cache.TTL(ctx, userKeySession)
		if err != nil {
			return response.ErrCodeInternalServer, nil, err
		}
		if state := cache.CheckTTL(ttl); state == cache.TTLHasValue {
			return response.ErrCodePendingVerification, &models.RegisterResponse{
				Token: hashEmail,
			}, nil
		}
		return response.ErrCodeEmailAlreadyExists, nil, nil
	}
	// Hash password
	hashedPassword, err := crypto.HashPassword(params.Password)
	if err != nil {
		return response.ErrCodeInternalServer, nil, err
	}

	// Generate OTP
	newOtp := random.GenerateSixDigitOtp()

	log.Info().Msgf(">>> OTP is: %d", newOtp)

	// Insert UserBase
	newUser, err := s.store.CreateUserBase(ctx, db.CreateUserBaseParams{
		UserEmail:    params.Email,
		UserPassword: hashedPassword,
		UserOtp:      strconv.Itoa(newOtp),
		UserHash:     hashEmail,
	})

	if err != nil {
		return response.ErrCodeInternalServer, nil, err
	}

	// Store otp to redis
	userKey := helpers.GetUserKeyOtp(hashEmail)
	err = s.cache.SetEx(ctx, userKey, newOtp, OTP_EXPIRATION)
	if err != nil {
		return response.ErrCodeInternalServer, nil, err
	}
	userKeySession := helpers.GetUserKeySession(hashEmail)
	err = s.cache.SetEx(ctx, userKeySession, newUser.UserID, OTP_EXPIRATION)
	if err != nil {
		return response.ErrCodeInternalServer, nil, err
	}

	// Send otp to email
	log.Info().Msg("Sending OTP to email")
	payload := &worker.PayloadSendVerifyEmail{
		Email: newUser.UserEmail,
		Otp:   strconv.Itoa(newOtp),
	}

	opts := []asynq.Option{
		asynq.MaxRetry(5),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.QueueCritical),
	}
	err = s.distributor.DistributeTaskSendVerifyEmail(ctx, payload, opts...)
	if err != nil {
		return response.ErrCodeInternalServer, nil, err
	}

	return response.ErrCodeSuccess, &models.RegisterResponse{
		Token: hashEmail,
	}, nil
}

func (s *authServiceImpl) GetTTLOtp(ctx context.Context, token string) (int, *models.VerifyOTPRes, error) {
	userKeyOtp := helpers.GetUserKeyOtp(token)
	ttl, err := s.cache.TTL(ctx, userKeyOtp)
	if err != nil {
		return response.ErrCodeInternalServer, nil, err
	}
	if ttl < 0 {
		return response.ErrCodeExpiredSession, nil, nil
	}
	return response.ErrCodeSuccess, &models.VerifyOTPRes{TTL: int(ttl.Seconds())}, nil
}

func (s *authServiceImpl) ResendOTP(ctx context.Context, params *models.ResendOTPReq) (int, error) {
	tokenFound, err := s.store.GetUserByUserHash(ctx, params.Token)
	if err != nil {
		if err == pgx.ErrNoRows {
			return response.ErrCodeInvalidParams, nil
		}
		return response.ErrCodeInternalServer, err
	}

	newOtp := random.GenerateSixDigitOtp()
	log.Info().Msgf(">>> OTP is: %d", newOtp)

	userKeyOtp := helpers.GetUserKeyOtp(params.Token)
	err = s.cache.SetEx(ctx, userKeyOtp, newOtp, OTP_EXPIRATION)
	if err != nil {
		return response.ErrCodeInternalServer, err
	}

	userKeySession := helpers.GetUserKeySession(params.Token)
	err = s.cache.SetEx(ctx, userKeySession, newOtp, OTP_EXPIRATION)
	if err != nil {
		return response.ErrCodeInternalServer, err
	}
	log.Info().Msg("Sending OTP to email")

	payload := &worker.PayloadSendVerifyEmail{
		Email: tokenFound.UserEmail,
		Otp:   strconv.Itoa(newOtp),
	}

	opts := []asynq.Option{
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.QueueCritical),
		asynq.MaxRetry(2),
	}

	err = s.distributor.DistributeTaskSendVerifyEmail(ctx, payload, opts...)
	if err != nil {
		return response.ErrCodeInternalServer, err
	}

	return response.ErrCodeSuccess, nil
}

// Logout implements services.IAuthService.
func (s *authServiceImpl) Logout(ctx context.Context) (code int, err error) {
	panic("unimplemented")
}

// RefreshToken implements services.IAuthService.
func (s *authServiceImpl) RefreshToken(ctx context.Context) (code int, res *models.LoginResponse, err error) {
	panic("unimplemented")
}

func NewAuthServiceImpl(store db.Store, cache cache.Cache, distributor worker.TaskDistributor) *authServiceImpl {
	return &authServiceImpl{
		store:       store,
		cache:       cache,
		distributor: distributor,
	}
}

var _ services.IAuthService = (*authServiceImpl)(nil)
