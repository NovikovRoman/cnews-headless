package routes

import (
	"github.com/NovikovRoman/cnews-headless/handlers"
	"github.com/NovikovRoman/cnews-headless/handlers/hlshell"
	"github.com/gofiber/fiber/v2"
)

func New(app *fiber.App, hlShell *hlshell.HeadlessShell) {
	app.Get("/", handlers.Homepage())

	app.Get("/html/", handlers.Html(hlShell))
	app.Get("/file/", handlers.File(hlShell))
}
