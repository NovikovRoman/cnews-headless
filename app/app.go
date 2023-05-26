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
	ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithLogf(log.Printf))
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
