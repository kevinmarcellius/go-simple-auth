package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dotenv-org/godotenvvault"
)

func main() {
	err := godotenvvault.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	hello := os.Getenv("HELLO")
	output := "Hello " + hello

	fmt.Printf(output)
}
