package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"news-app-api/internal/dto"
	"news-app-api/internal/usecase"
	"time"
)

const mediaSessionCookie = "media_session"

type MediaController struct {
	mediaUC usecase.MediaUseCase
}

func NewMediaController(mediaUC usecase.MediaUseCase) *MediaController {
	return &MediaController{mediaUC}
}

func (c *MediaController) Register() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var p dto.RegisterMediaParams
		if err := ctx.BodyParser(&p); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(newErrResponse(err))
		}

		media, err := c.mediaUC.Register(ctx.Context(), p)
		if err != nil {
			return err
		}

		ctx.Cookie(&fiber.Cookie{
			Name:     mediaSessionCookie,
			Value:    fmt.Sprint(media.ID),
			SameSite: "lax",
			Expires:  time.Now().Add(time.Hour * 24 * 30),
			HTTPOnly: true,
		})

		return ctx.Status(fiber.StatusCreated).JSON(newResponse(media))
	}
}

func (c *MediaController) Login() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var p dto.LoginMediaParams
		if err := ctx.BodyParser(&p); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(newErrResponse(err))
		}

		media, err := c.mediaUC.LoginMedia(ctx.Context(), p)
		if err != nil {
			return err
		}

		ctx.Cookie(&fiber.Cookie{
			Name:     mediaSessionCookie,
			Value:    fmt.Sprint(media.ID),
			SameSite: "lax",
			Expires:  time.Now().Add(time.Hour * 24 * 30),
			HTTPOnly: true,
		})

		return ctx.Status(fiber.StatusOK).JSON(newResponse(media))
	}
}

func (c *MediaController) Logout() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.ClearCookie(mediaSessionCookie)
		return ctx.SendStatus(fiber.StatusNoContent)
	}
}

func (c *MediaController) Authenticate() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		mediaID := ctx.Locals(mediaIDKey).(int64)

		media, err := c.mediaUC.GetMediaByID(ctx.Context(), mediaID)
		if err != nil {
			return err
		}

		ctx.Cookie(&fiber.Cookie{
			Name:     mediaSessionCookie,
			Value:    fmt.Sprint(media.ID),
			SameSite: "lax",
			Expires:  time.Now().Add(time.Hour * 24 * 30),
			HTTPOnly: true,
		})

		return ctx.Status(fiber.StatusOK).JSON(newResponse(media))
	}
}

func (c *MediaController) GetMediaList() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var p dto.GetMediaListParams
		if err := ctx.QueryParser(&p); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(newErrResponse(err))
		}

		res, err := c.mediaUC.GetMediaList(ctx.Context(), p)
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(newResponse(res))
	}
}

func (c *MediaController) ToggleSubscription() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var p dto.ToggleSubscriptionParams
		if err := ctx.ParamsParser(&p); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(newErrResponse(err))
		}

		p.UserID = ctx.Locals("userID").(int64)

		res, err := c.mediaUC.ToggleSubscription(ctx.Context(), p)
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(newResponse(res))
	}
}

func (c *MediaController) RegisterRoutes(r fiber.Router, mw *Middleware) {
	r.Post("/register", c.Register())
	r.Post("/login", c.Login())
	r.Post("/logout", c.Logout())
	r.Post("/authenticate", mw.AuthedMedia(), c.Authenticate())
	r.Get("", c.GetMediaList())
	r.Post("/:media_id/toggle-subscription", mw.AuthedUser(), c.ToggleSubscription())
}
