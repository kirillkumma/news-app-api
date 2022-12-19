package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"news-app-api/internal/dto"
	"news-app-api/internal/usecase"
)

type NewsController struct {
	newsUC usecase.NewsUseCase
}

func NewNewsController(newsUC usecase.NewsUseCase) *NewsController {
	return &NewsController{newsUC}
}

func (c *NewsController) CreateNews() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var p dto.CreateNewsParams
		if err := ctx.BodyParser(&p); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(newErrResponse(err))
		}

		p.MediaID = ctx.Locals(mediaIDKey).(int64)

		res, err := c.newsUC.CreateNews(ctx.Context(), p)
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusCreated).JSON(newResponse(res))
	}
}

func (c *NewsController) CreateOrUpdateAudio() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var p dto.CreateOrUpdateAudioParams
		if err := ctx.ParamsParser(&p); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(newErrResponse(err))
		}

		p.MediaID = ctx.Locals(mediaIDKey).(int64)

		f, err := ctx.FormFile("file")
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(newErrResponse(err))
		}

		p.File, err = f.Open()
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(newErrResponse(err))
		}
		defer p.File.Close()

		err = c.newsUC.CreateOrUpdateAudio(ctx.Context(), p)
		if err != nil {
			return err
		}

		return ctx.SendStatus(fiber.StatusNoContent)
	}
}

func (c *NewsController) GetAudio() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var p dto.GetAudioParams
		if err := ctx.ParamsParser(&p); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(newErrResponse(err))
		}

		audio, err := c.newsUC.GetAudio(ctx.Context(), p)
		if err != nil {
			return err
		}

		ctx.Set("content-length", fmt.Sprint(len(audio)))
		ctx.Set("content-type", "application/x-wav")
		return ctx.Status(fiber.StatusOK).Send(audio)
	}
}

func (c *NewsController) GetNews() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var p dto.GetNewsParams
		if err := ctx.ParamsParser(&p); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(newErrResponse(err))
		}

		news, err := c.newsUC.GetNews(ctx.Context(), p)
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(newResponse(news))
	}
}

func (c *NewsController) CreateOrUpdateImage() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var p dto.CreateOrUpdateImageParams
		if err := ctx.ParamsParser(&p); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(newErrResponse(err))
		}

		f, err := ctx.FormFile("file")
		if err != nil {
			return err
		}

		p.File, err = f.Open()
		if err != nil {
			return err
		}
		defer p.File.Close()

		p.MediaID = ctx.Locals(mediaIDKey).(int64)

		err = c.newsUC.CreateOrUpdateImage(ctx.Context(), p)
		if err != nil {
			return err
		}

		return ctx.SendStatus(fiber.StatusNoContent)
	}
}

func (c *NewsController) GetImage() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var p dto.GetImageParams
		if err := ctx.ParamsParser(&p); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(newErrResponse(err))
		}

		data, err := c.newsUC.GetImage(ctx.Context(), p)
		if err != nil {
			return err
		}

		ctx.Set("content-length", fmt.Sprint(len(data)))
		ctx.Set("content-type", "image/png")
		return ctx.Status(fiber.StatusOK).Send(data)
	}
}

func (c *NewsController) RegisterRoutes(r fiber.Router, mw *Middleware) {
	r.Post("", mw.AuthedMedia(), c.CreateNews())
	r.Put(":news_id/audio", mw.AuthedMedia(), c.CreateOrUpdateAudio())
	r.Get(":news_id/audio", c.GetAudio())
	r.Get(":news_id", c.GetNews())
	r.Put(":news_id/image", mw.AuthedMedia(), c.CreateOrUpdateImage())
	r.Get(":news_id/image", c.GetImage())
}
