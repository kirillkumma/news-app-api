package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"news-app-api/internal/dto"
	"news-app-api/internal/usecase"
	"time"
)

const userSessionCookie = "user_session"

type UserController struct {
	userUC usecase.UserUseCase
}

func NewUserController(userUC usecase.UserUseCase) *UserController {
	return &UserController{userUC}
}

func (c *UserController) Register() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var p dto.RegisterUserParams
		if err := ctx.BodyParser(&p); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(newErrResponse(err))
		}

		user, err := c.userUC.RegisterUser(ctx.Context(), p)
		if err != nil {
			return err
		}

		ctx.Cookie(&fiber.Cookie{
			Name:     userSessionCookie,
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
		var p dto.LoginUserParams
		if err := ctx.BodyParser(&p); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(newErrResponse(err))
		}

		user, err := c.userUC.LoginUser(ctx.Context(), p)
		if err != nil {
			return err
		}

		ctx.Cookie(&fiber.Cookie{
			Name:     userSessionCookie,
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
		ctx.ClearCookie(userSessionCookie)
		return ctx.SendStatus(fiber.StatusNoContent)
	}
}

func (c *UserController) Authenticate() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userID := ctx.Locals(userIDKey).(int64)

		user, err := c.userUC.GetUserByID(ctx.Context(), userID)
		if err != nil {
			return err
		}

		ctx.Cookie(&fiber.Cookie{
			Name:     userSessionCookie,
			Value:    fmt.Sprint(user.ID),
			HTTPOnly: true,
			Expires:  time.Now().Add(time.Hour * 24 * 30),
			SameSite: "lax",
		})

		return ctx.Status(fiber.StatusOK).JSON(newResponse(user))
	}
}

func (c *UserController) GetSubscriptionList() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var p dto.GetSubscriptionListParams
		if err := ctx.ParamsParser(&p); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(newErrResponse(err))
		}
		if err := ctx.QueryParser(&p); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(newErrResponse(err))
		}

		res, err := c.userUC.GetSubscriptionList(ctx.Context(), p)
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(newResponse(res))
	}
}

func (c *UserController) RegisterRoutes(r fiber.Router, mw *Middleware) {
	r.Post("register", c.Register())
	r.Post("login", c.Login())
	r.Post("logout", c.Logout())
	r.Post("authenticate", mw.AuthedUser(), c.Authenticate())
	r.Get(":user_id/subscriptions", c.GetSubscriptionList())
}
