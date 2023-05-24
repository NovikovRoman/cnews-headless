package handlers

import "github.com/gofiber/fiber/v2"

func Homepage() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		return c.Status(fiber.StatusOK).Send([]byte("(⌐■_■)"))
	}
}
