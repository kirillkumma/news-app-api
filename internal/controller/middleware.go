package controller

import (
	"github.com/gofiber/fiber/v2"
	"news-app-api/internal/dto"
	"strconv"
)

const userIDKey = "userID"

type Middleware struct {
}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) Auth() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userIDStr := ctx.Cookies("session")
		if userIDStr == "" {
			return &dto.AppError{
				Message: "Для совершения данной операции требуется авторизация",
				Code:    dto.ErrCodeUnauthorized,
			}
		}

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			return err
		}

		ctx.Locals(userIDKey, int64(userID))

		return ctx.Next()
	}
}
