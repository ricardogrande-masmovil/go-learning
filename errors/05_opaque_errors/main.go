package main

import (
	"errors"
	"fmt"
)

// Internal implementation details we don't want to leak.
// Raw DB or filesystem errors can leak SQL paths, library versions or tokens.
var errDBTimeout = errors.New("sql: connect postgres://user:sekret@localhost:5432: i/o timeout")

// fetchFromDatabase simulates a raw DB error occurring.
func fetchFromDatabase() error {
	return errDBTimeout
}

// GetUserLeaky leaks sensitive internal core details.
func GetUserLeaky() error {
	err := fetchFromDatabase()
	if err != nil {
		// BAD: Using %w or %v wraps/includes the sensitive DB error in the response
		// that the user/client will see.
		return fmt.Errorf("fetching user: %v", err)
	}
	return nil
}

// GetUserOpaque returns an opaque error to avoid leaking the kitchen core.
func GetUserOpaque() error {
	err := fetchFromDatabase()
	if err != nil {
		// GOOD: We log the internal error for our own diagnostics/monitoring...
		// log.Printf("internal db error: %v", err)

		// ... but we return an OPAQUE, safe error to the client.
		// Note there is no %w or %v here preserving the raw error value.
		// Context is still provided ("fetching user"), but the internal core is blocked.
		return errors.New("fetching user: internal service error")
	}
	return nil
}

func main() {
	fmt.Println("=== 05: Don't leak the kitchen core ===")

	leakyErr := GetUserLeaky()
	fmt.Println("Leaky Response   :", leakyErr)
	// Output: fetching user: sql: connect postgres://user:sekret@localhost:5432: i/o timeout
	// (Oops, we just leaked our DB connection string and credentials!)

	opaqueErr := GetUserOpaque()
	fmt.Println("Opaque Response  :", opaqueErr)
	// Output: fetching user: internal service error
	// (Safe and descriptive for external consumers)
}
