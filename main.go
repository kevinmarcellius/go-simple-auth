package main

import (
	"fmt"
	"log"
	
	"github.com/kevinmarcellius/go-simple-auth/config"

)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	hello := cfg.Postgres.Host
	
    output := "Hello " + hello


	fmt.Printf(output)
}
