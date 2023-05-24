package main

import (
	"github.com/NovikovRoman/cnews-headless/routes"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

func main() {
	var err error

	app := fiber.New()
	routes.New(app)
	if err = app.Listen("0.0.0.0:4444"); err != nil {
		log.Fatalf("error in ListenAndServe: %v", err)
	}
}
