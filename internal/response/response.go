package response

import "github.com/gofiber/fiber/v2"

type response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SuccessResponse(c *fiber.Ctx, code int, data interface{}) error {
	return c.JSON(response{
		Code:    code,
		Message: msg[code],
		Data:    data,
	})
}

func ErrorResponse(c *fiber.Ctx, code int, err error) error {
	return c.JSON(response{
		Code: code,
		Message: msg[code],
		Data: err.Error(),
	})
}
