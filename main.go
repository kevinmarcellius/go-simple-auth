package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/kevinmarcellius/go-simple-auth/config"
	handler "github.com/kevinmarcellius/go-simple-auth/internal/handler"

	repository "github.com/kevinmarcellius/go-simple-auth/internal/repository"
	service "github.com/kevinmarcellius/go-simple-auth/internal/service"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	hello := cfg.Postgres.Host

	output := "Hello " + hello

	fmt.Printf(output)

	db, err := config.ConnectPostgres(cfg.Postgres)
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}

	err = config.DBHealthCheck(db)
	if err != nil {
		log.Fatalf("Database health check failed: %v", err)
	}
	log.Println("Database connection is healthy.")

	healthHandler := handler.NewHealthHandler(db)

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository, cfg.JWTkey)
	userHandler := handler.NewUserHandler(userService)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, output)
	})

	v1 := e.Group("/api/v1")
	v1.GET("/health/ready", healthHandler.ReadinessCheck)
	v1.GET("/health/live", healthHandler.LivenessCheck)

	// user
	v1.POST("/user", userHandler.CreateUser)
	v1.POST("/user/login", userHandler.Login)

	// Start server

	port := fmt.Sprintf(":%d", cfg.Port)
	fmt.Printf("Starting server on port %s\n", port)
	e.Logger.Fatal(e.Start(port))
}
