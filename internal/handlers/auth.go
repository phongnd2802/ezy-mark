package handlers

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/phongnd2802/daily-social/internal/cache"
	"github.com/phongnd2802/daily-social/internal/db"
	"github.com/phongnd2802/daily-social/internal/dtos"
	"github.com/phongnd2802/daily-social/pkg/crypto"
	"github.com/phongnd2802/daily-social/pkg/random"
	"github.com/phongnd2802/daily-social/views/pages"
	"github.com/phongnd2802/daily-social/views/pages/auth"
	"github.com/redis/go-redis/v9"
)

const (
	OTP_EXPIRATION = 60
)

func (h *Handler) HandleLogin(c echo.Context) error {
	method := c.Request().Method
	if method == echo.GET {
		return render(c, auth.SignIn())
	}

	params := new(dtos.LoginRequest)
	if err := c.Bind(params); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	slog.Info("Params", "email", params.Email, "password", params.Password)

	return render(c, auth.SignIn())
}

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

			log.Println(">>>> OTP is:", newOtp)

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
			log.Println("Sent OTP to email")

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

	log.Println(">>>> OTP is:", newOtp)

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

	fmt.Println(newUser.CreatedAt)

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
	log.Println("Sent OTP to email")

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
			TTL: int(ttl.Seconds()), 
			Token: token,
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
