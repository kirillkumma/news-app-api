package controller

import "github.com/gofiber/fiber/v2"

func newResponse(data any) fiber.Map {
	return fiber.Map{
		"data": data,
	}
}

func newErrResponse(err error) fiber.Map {
	return fiber.Map{
		"error": fiber.Map{
			"message": err.Error(),
		},
	}
}
