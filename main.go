package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	
	"github.com/kevinmarcellius/go-simple-auth/config"

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
	fmt.Println("Database connection is healthy.")

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, output)
	})

	port:= fmt.Sprintf(":%d", cfg.Port)
	fmt.Printf("Starting server on port %s\n", port)
	e.Logger.Fatal(e.Start(port))
}
