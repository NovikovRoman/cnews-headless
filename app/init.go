package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
	if os.Getenv("DEBUG") == "" || os.Getenv("DEBUG") == "0" || os.Getenv("DEBUG") == "no" {
		log.SetLevel(log.WarnLevel)

	} else {
		log.SetLevel(log.InfoLevel)
	}
}
