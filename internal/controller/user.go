package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"news-app-api/internal/dto"
	"news-app-api/internal/usecase"
	"time"
)

type UserController struct {
	userUC usecase.UserUseCase
}

func NewUserController(userUC usecase.UserUseCase) *UserController {
	return &UserController{userUC}
}

func (c *UserController) Register() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var p dto.RegisterParams
		if err := ctx.BodyParser(&p); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(newErrResponse(err))
		}

		user, err := c.userUC.Register(ctx.Context(), p)
		if err != nil {
			return err
		}

		ctx.Cookie(&fiber.Cookie{
			Name:     "session",
			Value:    fmt.Sprint(user.ID),
			HTTPOnly: true,
			Expires:  time.Now().Add(time.Hour * 24 * 30),
			SameSite: "lax",
		})

		return ctx.Status(fiber.StatusCreated).JSON(newResponse(user))
	}
}

func (c *UserController) Login() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var p dto.LoginParams
		if err := ctx.BodyParser(&p); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(newErrResponse(err))
		}

		user, err := c.userUC.Login(ctx.Context(), p)
		if err != nil {
			return err
		}

		ctx.Cookie(&fiber.Cookie{
			Name:     "session",
			Value:    fmt.Sprint(user.ID),
			HTTPOnly: true,
			Expires:  time.Now().Add(time.Hour * 24 * 30),
			SameSite: "lax",
		})

		return ctx.Status(fiber.StatusOK).JSON(newResponse(user))
	}
}

func (c *UserController) Logout() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.ClearCookie("session")
		return ctx.SendStatus(fiber.StatusNoContent)
	}
}

func (c *UserController) Authenticate() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userID := ctx.Locals(userIDKey).(int64)

		user, err := c.userUC.GetByID(ctx.Context(), userID)
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(newResponse(user))
	}
}

func (c *UserController) RegisterRoutes(r fiber.Router, mw *Middleware) {
	r.Post("/register", c.Register())
	r.Post("/login", c.Login())
	r.Post("/logout", c.Logout())
	r.Post("/authenticate", mw.Auth(), c.Authenticate())
}
