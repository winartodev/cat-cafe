package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/winartodev/cat-cafe/internal/config"
	"github.com/winartodev/cat-cafe/internal/handlers"
	"github.com/winartodev/cat-cafe/internal/middleware"
	"github.com/winartodev/cat-cafe/internal/repositories"
	"github.com/winartodev/cat-cafe/internal/usecase"
	"github.com/winartodev/cat-cafe/pkg/jwt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	db, err := cfg.Database.SetupConnection()
	if err != nil {
		log.Fatalf("Could setup database: %v", err)
	}

	redisClient, err := cfg.Redis.SetupRedisClient()
	if err != nil {
		log.Fatalf("Could setup redis: %v", err)
	}

	app := fiber.New(fiber.Config{
		AppName: cfg.App.Name,
	})
	app.Use(logger.New())
	app.Use(cors.New())

	jwtManager := jwt.NewJWT(cfg.JWT.SecretKey, cfg.JWT.TokenDuration)

	repo := repositories.SetupRepository(db, redisClient)

	uc := usecase.SetUpUseCase(*repo, jwtManager)

	middleware_ := middleware.NewMiddleware(jwtManager, repo.UserRepository)
	handlers.SetupHandler(app, *uc, middleware_)

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
