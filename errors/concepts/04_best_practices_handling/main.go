package main

import (
	"fmt"
	"log"
	"os"
)

// CallService simulates an external call.
func CallService() error {
	return fmt.Errorf("connection refused")
}

// FetchUserBad demonstrates bad practice.
// Handle or return. Never both!
func FetchUserBad() error {
	err := CallService()
	if err != nil {
		// BAD: We are logging the error here...
		log.Printf("call service foo: %v", err)
		// ... AND returning it, causing the caller to likely log it again!
		return fmt.Errorf("call service foo: %w", err)
	}
	return nil
}

// FetchUserGood demonstrates good practice.
// Delegate logging to client. Avoid repetition.
func FetchUserGood() error {
	err := CallService()
	if err != nil {
		// GOOD: Just add context and return. Leave logging to the top level.
		return fmt.Errorf("call service foo: %w", err)
	}
	return nil
}

// FetchUserGraceful demonstrates graceful degradation.
// If the error is recoverable or non-critical, log it and proceed without breaking.
func FetchUserGraceful() error {
	err := CallService()
	if err != nil {
		// Example: The external service is just a cache. If it fails, we
		// gracefully degrade by logging the issue and skipping it instead of failing.
		log.Printf("cache service unavailable, degrading gracefully: %v", err)
		return nil
	}
	return nil
}

func main() {
	// Redirect log output to avoid timestamps changing output unpredictably
	// for the sake of clear demonstration.
	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	fmt.Println("=== 04: Handle or Return. Never both. ===")

	fmt.Println("\n--- BAD PRACTICE (Logging and Returning) ---")
	err := FetchUserBad()
	if err != nil {
		log.Printf("main process: fetch user: %v", err)
	}
	// Output has annoying duplicates:
	// call service foo: connection refused
	// main process: fetch user: call service foo: connection refused

	fmt.Println("\n--- GOOD PRACTICE (Delegate logging to client) ---")
	err = FetchUserGood()
	if err != nil {
		log.Printf("main process: fetch user: %v", err)
	}
	// Output is clean and linear:
	// main process: fetch user: call service foo: connection refused

	fmt.Println("\n--- GOOD PRACTICE (Match and Degrade Gracefully) ---")
	err = FetchUserGraceful()
	if err != nil {
		log.Printf("main process: fetch user: %v", err)
	}
	// Output shows the cache missed but execution continued cleanly.
}
