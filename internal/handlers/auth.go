package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/phongnd2802/daily-social/internal/dtos"
	"github.com/phongnd2802/daily-social/views/pages/auth"
)

func HandleLogin(c echo.Context) error {
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

func HandleRegister(c echo.Context) error {
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
		errMsg["password"] = params.Password
		errMsg["errPassword"] = "Password must be at least 8 characters long."
		return render(c, auth.SignUp(errMsg))
	}

	if params.Password != params.ConfirmPassword {
		errMsg["password"] = params.Password
		errMsg["match"] = params.ConfirmPassword
		errMsg["errMatch"] = "Passwords do not match, Please try again."
		return render(c, auth.SignUp(errMsg))
	}



	return render(c, auth.SignUp(nil))
}


func HandleVerifyOTP(c echo.Context) error {
	method := c.Request().Method
	if method == echo.GET {
		return render(c, auth.VerifyOTP())
	}

	return nil
}