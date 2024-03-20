package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {

	godotenv.Load()
	portListener := os.Getenv("PORT")
	if portListener == "" {
		log.Fatalf("PORT is not set")
	}
	fmt.Println("Server is running on port: ", portListener)
}
