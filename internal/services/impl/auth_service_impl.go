package impl

import (
	"context"
	"strconv"
	"time"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5"
	"github.com/phongnd2802/daily-social/internal/cache"
	"github.com/phongnd2802/daily-social/internal/consts"
	"github.com/phongnd2802/daily-social/internal/db"
	"github.com/phongnd2802/daily-social/internal/dtos"
	"github.com/phongnd2802/daily-social/internal/helpers"
	"github.com/phongnd2802/daily-social/internal/response"
	"github.com/phongnd2802/daily-social/internal/worker"
	"github.com/phongnd2802/daily-social/pkg/crypto"
	"github.com/phongnd2802/daily-social/pkg/random"
	"github.com/rs/zerolog/log"
)

type authServiceImpl struct {
	store       db.Store
	cache       cache.Cache
	distributor worker.TaskDistributor
}

func (s *authServiceImpl) Login(ctx context.Context, params *dtos.LoginRequest) (code int, res *dtos.LoginResponse, err error) {

	return response.ErrCodeSuccess, nil, nil
}

func (s *authServiceImpl) VerifyOTP(ctx context.Context, params *dtos.VerifyOTPReq) (int, *dtos.VerifyOTPRes, error) {
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
		return response.ErrCodeOtpDoesNotMatch, &dtos.VerifyOTPRes{
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

func (s *authServiceImpl) Register(ctx context.Context, params *dtos.RegisterRequest) (int, *dtos.RegisterResponse, error) {
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
			return response.ErrCodePendingVerification, &dtos.RegisterResponse{
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
	err = s.cache.SetEx(ctx, userKey, newOtp, consts.OTP_EXPIRATION)
	if err != nil {
		return response.ErrCodeInternalServer, nil, err
	}
	userKeySession := helpers.GetUserKeySession(hashEmail)
	err = s.cache.SetEx(ctx, userKeySession, newUser.UserID, consts.OTP_EXPIRATION)
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

	return response.ErrCodeSuccess, &dtos.RegisterResponse{
		Token: hashEmail,
	}, nil
}

func (s *authServiceImpl) GetTTLOtp(ctx context.Context, token string) (int, *dtos.VerifyOTPRes, error) {
	userKeyOtp := helpers.GetUserKeyOtp(token)
	ttl, err := s.cache.TTL(ctx, userKeyOtp)
	if err != nil {
		return response.ErrCodeInternalServer, nil, err
	}
	if ttl < 0 {
		return response.ErrCodeExpiredSession, nil, nil
	}
	return response.ErrCodeSuccess, &dtos.VerifyOTPRes{TTL: int(ttl.Seconds())} ,nil
}


func NewAuthServiceImpl(store db.Store, cache cache.Cache, distributor worker.TaskDistributor) *authServiceImpl {
	return &authServiceImpl{
		store:       store,
		cache:       cache,
		distributor: distributor,
	}
}
