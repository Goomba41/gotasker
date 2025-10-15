package main

import (
	// "os"

	"log"

	// "github.com/gin-gonic/gin"

	"Goomba41/gotasker/pkg/configuration"
	"Goomba41/gotasker/pkg/database"
)

func main() {
	cfg, err := configuration.Init()
	if err != nil {
		log.Fatalf("Configuration file error: %v", err)
	}

	if err := database.SetConfig(cfg.Database); err != nil {
		log.Fatalf("Database DSN set error: %v", err)
	}

	_, err = database.Connect()
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
}
