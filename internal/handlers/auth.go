package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/phongnd2802/daily-social/internal/cache"
	"github.com/phongnd2802/daily-social/internal/db"
	"github.com/phongnd2802/daily-social/internal/dtos"
	"github.com/phongnd2802/daily-social/internal/worker"
	"github.com/phongnd2802/daily-social/pkg/crypto"
	"github.com/phongnd2802/daily-social/pkg/random"
	"github.com/phongnd2802/daily-social/pkg/token"
	"github.com/phongnd2802/daily-social/pkg/utils"
	"github.com/phongnd2802/daily-social/views/pages"
	"github.com/phongnd2802/daily-social/views/pages/auth"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

const (
	OTP_EXPIRATION = 60
)

func (h *Handler) HandleLogin(c echo.Context) error {
	method := c.Request().Method
	if method == echo.GET {
		return render(c, auth.SignIn(auth.SignInViewProps{}))
	}

	params := new(dtos.LoginRequest)
	if err := c.Bind(params); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	errs := make(map[string]string)
	// Check Email
	userFound, err := h.store.GetUserBaseByEmail(c.Request().Context(), params.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			errs["errLogin"] = "The email or password is incorrect!"
			return render(c, auth.SignIn(auth.SignInViewProps{
				Errors: errs,
			}))
		}
		return err
	}
	match := crypto.VerifyPassword(params.Password, userFound.UserPassword)
	if !match {
		errs["errLogin"] = "The email or password is incorrect!"
		return render(c, auth.SignIn(auth.SignInViewProps{
			Errors: errs,
		}))
	}

	// Generate Token
	subToken := utils.GenerateCliTokenUUID(userFound.UserID)
	log.Info().Str("subToken", subToken).Msg("Generate Token")

	accessToken, err := token.CreateToken(subToken, "1h")
	if err != nil {
		return err
	}
	
	refreshToken, err := token.CreateToken(subToken, "168h")
	if err != nil {
		return err
	}

	c.SetCookie(&http.Cookie{
		Name: "access-token",
		Value: accessToken,
		HttpOnly: true,
		Secure: false,
		SameSite: http.SameSiteLaxMode,
		Path: "/",
		Expires: time.Now().Add(1 * time.Hour),
	})

	c.SetCookie(&http.Cookie{
		Name: "refresh-token",
		Value: refreshToken,
		HttpOnly: true,
		Secure: false,
		SameSite: http.SameSiteLaxMode,
		Path: "/",
		Expires: time.Now().Add(7 * 24 * time.Hour),
	})

	return c.Redirect(http.StatusSeeOther, "/")
}

// /////////////////////////////////////////
// /										///
// /				Register			   ///
// /									  ///
// /////////////////////////////////////////
func (h *Handler) HandleRegister(c echo.Context) error {
	method := c.Request().Method
	if method == echo.GET {
		return render(c, auth.SignUp(nil))
	}

	params := new(dtos.RegisterRequest)
	errMsg := make(map[string]string)
	if err := c.Bind(params); err != nil {
		errMsg["errParams"] = "Invalid request, check your input"
		return render(c, auth.SignUp(errMsg))
	}

	fmt.Println("Email: ", params.Email)
	fmt.Println("Password: ", params.Password)

	if len(params.Password) < 8 {
		errMsg["email"] = params.Email
		errMsg["password"] = params.Password
		errMsg["errPassword"] = "Password must be at least 8 characters long."
		return render(c, auth.SignUp(errMsg))
	}

	if params.Password != params.ConfirmPassword {
		errMsg["email"] = params.Email
		errMsg["password"] = params.Password
		errMsg["match"] = params.ConfirmPassword
		errMsg["errMatch"] = "Passwords do not match, Please try again."
		return render(c, auth.SignUp(errMsg))
	}

	// Hash Email
	hashUserEmail := crypto.GetHash(params.Email)
	// Check Email Exists
	userFound, err := h.store.GetUserBaseByEmail(c.Request().Context(), params.Email)
	if err != nil && err != pgx.ErrNoRows {
		return err
	}
	if userFound.UserID > 0 {

		userKeySession := getUserKeySession(hashUserEmail)
		ttl, err := h.cache.TTL(c.Request().Context(), userKeySession)
		if err != nil {
			return err
		}
		if state := cache.CheckTTL(ttl); state == cache.TTLHasValue {
			return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/verify-otp?token=%s", hashUserEmail))
		}

		if !userFound.IsVerified.Bool {
			// Generate OTP
			newOtp := random.GenerateSixDigitOtp()

			log.Debug().Msgf(">>>> OTP is: %d", newOtp)

			userKey := getUserKeyOtp(hashUserEmail)
			err = h.cache.SetEx(c.Request().Context(), userKey, newOtp, OTP_EXPIRATION)
			if err != nil {
				return err
			}

			userKeySession := getUserKeySession(hashUserEmail)
			err = h.cache.SetEx(c.Request().Context(), userKeySession, userFound.UserID, OTP_EXPIRATION)
			if err != nil {
				return err
			}
			// Send Otp to email
			log.Info().Msg("Sending OTP to email")

			return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/verify-otp?token=%s", hashUserEmail))
		}
		errMsg["email"] = params.Email
		errMsg["errEmailExists"] = "Email already exists."
		return render(c, auth.SignUp(errMsg))
	}

	// Hash password
	hashedPassword, err := crypto.HashPassword(params.Password)
	if err != nil {
		return render(c, auth.SignUp(nil))
	}

	fmt.Println(hashedPassword)

	// Generate OTP
	newOtp := random.GenerateSixDigitOtp()

	log.Info().Msgf(">>>> OTP is: %d", newOtp)

	// Insert UserBase
	newUser, err := h.store.CreateUserBase(c.Request().Context(), db.CreateUserBaseParams{
		UserEmail:    params.Email,
		UserPassword: hashedPassword,
		UserOtp:      strconv.Itoa(newOtp),
		UserHash:     hashUserEmail,
	})
	if err != nil {
		return err
	}

	// Store Otp to redis
	userKey := getUserKeyOtp(hashUserEmail)
	err = h.cache.SetEx(c.Request().Context(), userKey, newOtp, OTP_EXPIRATION)
	if err != nil {
		return err
	}

	userKeySession := getUserKeySession(hashUserEmail)
	err = h.cache.SetEx(c.Request().Context(), userKeySession, newUser.UserID, OTP_EXPIRATION)
	if err != nil {
		return err
	}
	// Send Otp to email
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
	err = h.distributor.DistributeTaskSendVerifyEmail(c.Request().Context(), payload, opts...)
	if err != nil {
		return err
	}
	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/verify-otp?token=%s", hashUserEmail))
}

func (h *Handler) HandleVerifyOTP(c echo.Context) error {
	method := c.Request().Method
	if method == echo.GET {
		token := c.QueryParam("token")
		if token == "" {
			return c.Redirect(http.StatusSeeOther, "/signin")
		}

		userKeySession := getUserKeySession(token)
		ttl, err := h.cache.TTL(c.Request().Context(), userKeySession)
		if err != nil {
			return err
		}
		if state := cache.CheckTTL(ttl); state == cache.TTLExpired {
			return c.Redirect(http.StatusSeeOther, "/signin")
		}
		sess, _ := session.Get("session", c)
		errMsg, _ := sess.Values["error"].(string)
		delete(sess.Values, "error")
		sess.Save(c.Request(), c.Response())

		return render(c, auth.VerifyOTP(auth.VerifyOTPViewProps{
			TTL:    int(ttl.Seconds()),
			Token:  token,
			Errors: map[string]string{"errOtp": errMsg},
		}))
	}

	//errMsg := make(map[string]string)
	params := new(dtos.VerifyOTPReq)
	if err := c.Bind(params); err != nil {
		return err
	}
	fmt.Println(">>>> OTP is:", params.Otp)
	fmt.Println(">>> Token is:", params.Token)

	// Get OTP from redis
	userKeyOtp := getUserKeyOtp(params.Token)
	otpFound, err := h.cache.Get(c.Request().Context(), userKeyOtp)
	if err != nil {
		if err == redis.Nil {
			return render(c, auth.VerifyOTP(auth.VerifyOTPViewProps{}))
		}
		return err
	}

	if otpFound != params.Otp {
		sess, _ := session.Get("session", c)
		sess.Values["error"] = "OTP does not match!"
		sess.Save(c.Request(), c.Response())
		//errMsg["errOtp"] = "OTP does not match!"
		// return render(c, auth.VerifyOTP(auth.VerifyOTPViewProps{
		// 	Errors: errMsg,
		// 	Token:  params.Token,
		// 	TTL:    int(ttl.Seconds()),
		// }))
		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/verify-otp?token=%s", params.Token))
	}

	// Updating User Veifited
	userUpdated, err := h.store.UpdateUserVerify(c.Request().Context(), params.Token)
	if err != nil {
		return err
	}

	// Del Cache
	keys := []string{getUserKeyOtp(params.Token), getUserKeySession(params.Token)}
	err = h.cache.Del(c.Request().Context(), keys...)
	if err != nil {
		return err
	}

	// Create User Profile
	nickName := extractNameFromEmail(userUpdated.UserEmail)
	_, err = h.store.CreateUserProfile(c.Request().Context(), db.CreateUserProfileParams{
		UserID:       userUpdated.UserID,
		UserEmail:    userUpdated.UserEmail,
		UserNickname: nickName,
	})
	if err != nil {
		return err
	}
	return c.Redirect(http.StatusSeeOther, "/signin")
}

func (h *Handler) HandleForgotPassword(c echo.Context) error {
	return render(c, pages.Index())
}
