package controller

import (
	"github.com/gofiber/fiber/v2"
	"news-app-api/internal/dto"
	"news-app-api/internal/usecase"
)

type FeedController struct {
	feedUC usecase.FeedUseCase
}

func NewFeedController(feedUC usecase.FeedUseCase) *FeedController {
	return &FeedController{feedUC}
}

func (c *FeedController) GetFeed() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var p dto.GetFeedParams
		if err := ctx.QueryParser(&p); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(newErrResponse(err))
		}

		p.UserID = ctx.Locals(userIDKey).(int64)

		res, err := c.feedUC.GetFeed(ctx.Context(), p)
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(newResponse(res))
	}
}

func (c *FeedController) RegisterRoutes(r fiber.Router, mw *Middleware) {
	r.Get("", mw.AuthedUser(), c.GetFeed())
}
