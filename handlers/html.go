package handlers

import (
	"context"
	"encoding/json"
	"html"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

func Html() fiber.Handler {
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

		selector := c.Query("selector", "")
		if res.Html, err = getHtml(c.Context(), u, selector); err != nil {
			res.Error = err.Error()
			log.Errorf("%s %v", u, err)
		}
		return c.Status(fiber.StatusOK).JSON(res)
	}
}

func getHtml(ctx context.Context, target, selector string) (body string, err error) {
	opts := []chromedp.ExecAllocatorOption{
		// chromedp.ExecPath("/headless-shell"),
		//chromedp.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.5672.127 Safari/537.36"),
		//chromedp.WindowSize(1280, 720),
		chromedp.NoFirstRun,
		//chromedp.NoDefaultBrowserCheck,
		chromedp.Headless,
		chromedp.IgnoreCertErrors,
		chromedp.DisableGPU,
	}

	ctx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx, chromedp.WithLogf(log.Printf))
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 45*time.Second)
	defer cancel()

	if selector == "" {
		selector = "body"
	}
	start := time.Now()
	err = chromedp.Run(
		ctx,
		network.Enable(),
		chromedp.Navigate(target),
		chromedp.WaitReady(selector),
		// cookie(),
		chromedp.OuterHTML("html", &body, chromedp.ByQuery),
		// chromedp.WaitNotVisible(`#trk_jschal_nojs`, chromedp.ByQuery),
		// cookie(),
		//chromedp.FullScreenshot(&b, 80),
		//removeCookie(),
	)
	log.Infof("%s %f sec", target, time.Since(start).Seconds())
	if err != nil {
		return
	}

	body = html.UnescapeString(body)
	return
}
