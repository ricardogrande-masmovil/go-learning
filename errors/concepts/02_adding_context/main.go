package main

import (
	"errors"
	"fmt"
)

// ErrDatabaseLost simulates a raw infrastructure error.
var ErrDatabaseLost = errors.New("database connection lost")

// connectDB simulates a core raw error (e.g., DB failure).
func connectDB() error {
	return ErrDatabaseLost
}

// initStoreObfuscated demonstrates intentional obfuscation using %v.
// The %v verb hides the underlying error type/identity from callers,
// preventing them from matching against it via errors.Is/errors.As.
// This is useful when you don't want callers to rely on internal details.
func initStoreObfuscated() error {
	err := connectDB()
	if err != nil {
		// Obfuscate the root!
		// The %v verb creates a new error string and strips the original identity.
		return fmt.Errorf("creating store: %v", err)
	}
	return nil
}

// initStorePreserved demonstrates the RIGHT way: adding context, preserving the root.
// This is "Building the Go-nion".
func initStorePreserved() error {
	err := connectDB()
	if err != nil {
		// Context added, root preserved!
		// The %w verb wraps the error.
		// The outer layer is the added context ("creating store").
		// The inner layer (core) is the original raw error.
		return fmt.Errorf("creating store: %w", err)
	}
	return nil
}

func main() {
	fmt.Println("=== 02: Building the Go-nion (Adding layers of context) ===")

	err1 := initStoreObfuscated()
	fmt.Printf("Obfuscated Wrapping   : %v\n", err1)
	// Because we used %v, callers CANNOT match against the DB error.
	// Output will be false.
	fmt.Printf("Is ErrDatabaseLost?   : %v\n\n", errors.Is(err1, ErrDatabaseLost))

	err2 := initStorePreserved()
	fmt.Printf("Preserved Wrapping    : %v\n", err2)
	// Because we used %w, we can inspect and peel the layers!
	// Output will be true.
	fmt.Printf("Is ErrDatabaseLost?   : %v\n", errors.Is(err2, ErrDatabaseLost))
}
