package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/winartodev/cat-cafe/internal/config"
	"github.com/winartodev/cat-cafe/internal/controllers"
	"github.com/winartodev/cat-cafe/internal/handlers"
	"github.com/winartodev/cat-cafe/internal/repositories"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	cfg.Database.SetupConnection()

	app := fiber.New(fiber.Config{
		AppName: cfg.App.Name,
	})

	repo := repositories.SetupRepository()

	ctrl := controllers.SetUpController(*repo)

	handlers.SetupHandler(app, *ctrl)

	go func() {
		port := fmt.Sprintf(":%d", cfg.App.Port)
		if err := app.Listen(port); err != nil {
			log.Panicf("Failed to start server : %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	fmt.Println("Shutting down server....")
	_ = app.Shutdown()
}
