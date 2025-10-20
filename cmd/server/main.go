package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	// "github.com/gin-gonic/gin"

	repositories "goomba41/gotasker/internal/repository"
	"goomba41/gotasker/internal/repository/db"
	"goomba41/gotasker/internal/dto"
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

	connection, err := database.Connect()
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer connection.Close()

	queries := db.New(connection)

	userRepo := repositories.NewUserRepository(queries, connection)

	createdUser, err := userRepo.Create(context.Background(), "anton.borodawkin@yandex.ru", "password")
	if err == nil {
		log.Printf("Created user: %v", createdUser)
	}

	gettedUser, err := userRepo.GetByEmail(context.Background(), "anton.borodawkin@yandex.ru")
	if err == nil {
		log.Printf("Getted user: %v", gettedUser)
	}

	newPassword := fmt.Sprintf("%d.%s", gettedUser.ID, gettedUser.Password)
	patchData := dto.UserPatch{
		Password: &newPassword,
	}

	patchedUser, err := userRepo.Patch(context.Background(), 5, patchData)
	if err == nil {
		log.Printf("Patched user: %v", patchedUser)
	}
}
