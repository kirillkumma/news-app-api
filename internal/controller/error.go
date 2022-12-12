package controller

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"news-app-api/internal/dto"
)

func ErrHandler(ctx *fiber.Ctx, err error) error {
	var appErr *dto.AppError
	if errors.As(err, &appErr) {
		return ctx.Status(appErr.Code).JSON(newErrResponse(err))
	}
	log.Error(err.Error())
	return ctx.SendStatus(fiber.StatusInternalServerError)
}
