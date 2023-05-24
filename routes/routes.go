package routes

import (
	"github.com/NovikovRoman/cnews-headless/handlers"
	"github.com/gofiber/fiber/v2"
)

func New(app *fiber.App) {
	app.Get("/", handlers.Homepage())

	app.Get("/html/", handlers.Html())
}
