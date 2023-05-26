package handlers

import (
	"encoding/json"

	"github.com/NovikovRoman/cnews-headless/handlers/hlshell"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

func File(hlShell *hlshell.HeadlessShell) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		type result struct {
			Body  []byte `json:"body"`
			Error string `json:"error"`
		}

		res := result{}
		u := c.Query("url", "")
		if u == "" {
			res.Error = "url is empty"
			b, _ := json.Marshal(res)
			return c.Status(fiber.StatusOK).Send(b)
		}

		if res.Body, err = hlShell.File(u); err != nil {
			res.Error = err.Error()
			log.Errorf("%s %v", u, err)
		}
		return c.Status(fiber.StatusOK).JSON(res)
	}
}
