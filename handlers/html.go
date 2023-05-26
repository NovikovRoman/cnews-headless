package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/NovikovRoman/cnews-headless/handlers/webpage"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

func Html() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		type request struct {
			Url      string   `json:"url"`
			Selector string   `json:"selector"`
			Cookies  []string `json:"cookies"`
		}

		type result struct {
			Html    string   `json:"html"`
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
		if err = wp.Get(c.Context(), req.Selector, req.Cookies...); err != nil {
			res.Error = err.Error()
			log.Errorf("%s %v", req.Url, err)

		} else {
			res.Html = wp.String()
			res.Cookies = wp.Cookies()
		}
		return c.Status(fiber.StatusOK).JSON(res)
	}
}
