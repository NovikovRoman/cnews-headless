package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/NovikovRoman/cnews-headless/handlers/webpage"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

func File() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		type request struct {
			Url     string   `json:"url"`
			Cookies []string `json:"cookies"`
		}

		type result struct {
			Body    []byte   `json:"body"`
			Cookies []string `json:"cookies"`
			Error   string   `json:"error"`
		}

		res := result{}
		var req request
		if err = c.BodyParser(&req); err != nil {
			res.Error = fmt.Sprintf("Invalid request %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(res)
		}

		if req.Url == "" {
			res.Error = "url is empty"
			b, _ := json.Marshal(res)
			return c.Status(fiber.StatusOK).Send(b)
		}

		wp := webpage.New(req.Url)
		if err = wp.File(c.Context(), req.Cookies...); err != nil {
			res.Error = err.Error()
			log.Errorf("%s %v", req.Url, err)

		} else {
			res.Body = wp.Bytes()
			res.Cookies = wp.Cookies()
		}
		return c.Status(fiber.StatusOK).JSON(res)
	}
}
