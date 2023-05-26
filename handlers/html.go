package handlers

import (
	"encoding/json"

	"github.com/NovikovRoman/cnews-headless/handlers/hlshell"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

func Html(hlShell *hlshell.HeadlessShell) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		type result struct {
			Html  string `json:"html"`
			Error string `json:"error"`
		}

		res := result{}
		u := c.Query("url", "")
		if u == "" {
			res.Error = "url is empty"
			b, _ := json.Marshal(res)
			return c.Status(fiber.StatusOK).Send(b)
		}

		if res.Html, err = hlShell.Html(u, c.Query("selector", "")); err != nil {
			res.Error = err.Error()
			log.Errorf("%s %v", u, err)
		}
		return c.Status(fiber.StatusOK).JSON(res)
	}
}
