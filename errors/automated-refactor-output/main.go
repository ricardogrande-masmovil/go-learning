package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/ricardogrande-masmovil/go-learning/errors/automated-refactor-output/domain"
	"github.com/ricardogrande-masmovil/go-learning/errors/automated-refactor-output/service"
)

func main() {
	fmt.Println("Starting automated refactor test program...")

	user, err := service.GetUserProfile("12345")
	if err != nil {
		var domainErr *domain.UserDomainError
		if errors.As(err, &domainErr) {
			log.Printf("[SYSTEM ERR] User registration failed: %v", domainErr.Err)
			fmt.Printf("User Error (to client): %s\n", domainErr.UserMsg)
			return
		}

		log.Fatalf("Critical entrypoint failure: %v", err)
	}

	fmt.Printf("Successfully fetched user: %s (%s)\n", user.Name, user.ID)
}
