package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	// "github.com/gin-gonic/gin"

	"goomba41/gotasker/pkg/configuration"
	"goomba41/gotasker/pkg/database"
)

func main() {
	configPath := flag.String("config", "", "Configuration file path")
	flag.Parse()

	if *configPath == "" {
		fmt.Fprintf(flag.CommandLine.Output(), "Error: -config flag is required\n\n")
		flag.Usage()
		fmt.Println()
		os.Exit(1)
	}

	cfg, err := configuration.Init(*configPath)
	if err != nil {
		log.Fatalf("Configuration file: %v", err)
	}

	if err := database.SetConfig(cfg.Database); err != nil {
		log.Fatalf("Database DSN set error: %v", err)
	}

	_, err = database.Connect()
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
}
