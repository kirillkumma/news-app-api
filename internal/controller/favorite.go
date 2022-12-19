package controller

import (
	"github.com/gofiber/fiber/v2"
	"news-app-api/internal/dto"
	"news-app-api/internal/usecase"
)

type FavoriteController struct {
	newsUC usecase.NewsUseCase
}

func NewFavoriteController(newsUC usecase.NewsUseCase) *FavoriteController {
	return &FavoriteController{newsUC}
}

func (c *FavoriteController) GetFavoriteList() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var p dto.GetFavoriteListParams
		if err := ctx.QueryParser(&p); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(newErrResponse(err))
		}

		p.UserID = ctx.Locals(userIDKey).(int64)

		res, err := c.newsUC.GetFavoriteList(ctx.Context(), p)
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(newResponse(res))
	}
}

func (c *FavoriteController) RegisterRoutes(r fiber.Router, mw *Middleware) {
	r.Get("", mw.AuthedUser(), c.GetFavoriteList())
}
