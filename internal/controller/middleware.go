package controller

import (
	"github.com/gofiber/fiber/v2"
	"news-app-api/internal/dto"
	"strconv"
)

const userIDKey = "userID"
const mediaIDKey = "mediaID"

type Middleware struct {
}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) AuthedUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userIDStr := ctx.Cookies(userSessionCookie)
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

func (m *Middleware) AuthedMedia() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		mediaIDStr := ctx.Cookies(mediaSessionCookie)
		if mediaIDStr == "" {
			return &dto.AppError{
				Message: "Для совершения данной операции требуется авторизация",
				Code:    dto.ErrCodeUnauthorized,
			}
		}

		mediaID, err := strconv.Atoi(mediaIDStr)
		if err != nil {
			return err
		}

		ctx.Locals(mediaIDKey, int64(mediaID))

		return ctx.Next()
	}
}

func (m *Middleware) OptionalAuthedUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userIDStr := ctx.Cookies(userSessionCookie)
		if userIDStr == "" {
			return ctx.Next()
		}

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			return ctx.Next()
		}

		ctx.Locals(userIDKey, int64(userID))

		return ctx.Next()
	}
}
