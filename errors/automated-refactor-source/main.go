package main

import (
	"fmt"
	"log"

	"github.com/ricardogrande-masmovil/go-learning/errors/automated-refactor-source/service"
)

func main() {
	fmt.Println("Starting automated refactor test program...")

	user, err := service.GetUserProfile("12345")
	if err != nil {
		// Just logging and exiting. The inner codes have the anti-patterns.
		log.Fatalf("Critical entrypoint failure: %v", err)
	}

	fmt.Printf("Successfully fetched user: %s (%s)\n", user.Name, user.ID)
}
