package main

import (
	"context"

	"github.com/NovikovRoman/cnews-headless/handlers/hlshell"
	"github.com/NovikovRoman/cnews-headless/routes"
	"github.com/chromedp/chromedp"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

func main() {
	opts := []chromedp.ExecAllocatorOption{
		chromedp.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"),
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Headless,
		chromedp.IgnoreCertErrors,
		chromedp.DisableGPU,
	}

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx, chromedp.WithLogf(log.Printf))
	defer cancel()

	err := chromedp.Run(ctx)
	if err != nil {
		log.Fatalf("start headless-shell: %v", err)
	}

	app := fiber.New()
	routes.New(app, hlshell.New(ctx))
	if err = app.Listen("0.0.0.0:4444"); err != nil {
		log.Fatalf("error in ListenAndServe: %v", err)
	}
}
