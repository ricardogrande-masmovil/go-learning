package main

import (
	"errors"
	"fmt"
	"os"
)

// The single return value form of a type assertion will panic on an incorrect type.
// BAD: Always use the "comma ok" idiom to handle type assertion failures.
func getGreetingBad(i interface{}) string {
	// Panics if i is not a string
	s := i.(string)
	return fmt.Sprintf("Hello %s", s)
}

// GOOD: Use "comma ok" to fail gracefully.
func getGreetingGood(i interface{}) (string, error) {
	s, ok := i.(string)
	if !ok {
		return "", errors.New("input is not a string")
	}
	return fmt.Sprintf("Hello %s", s), nil
}

// runBad uses panics for control flow/validation.
// BAD: Panic/recover is not an error handling strategy.
func runBad(args []string) {
	if len(args) == 0 {
		panic("an argument is required") // Don't do this! Use for irrecoverable only.
	}
}

// runGood uses returned errors for control flow/validation.
// GOOD: Return errors and let the caller decide how to handle them.
func runGood(args []string) error {
	if len(args) == 0 {
		return errors.New("an argument is required")
	}
	return nil
}

func main() {
	fmt.Println("=== 07: Don't Panic ===")
	fmt.Println("Code running in production must avoid panics. Panics cause cascading failures.")

	// Example 1: Type Assertions
	res, err := getGreetingGood(123)
	if err != nil {
		fmt.Printf("Handled type error gracefully: %v\n", err)
	} else {
		fmt.Println(res)
	}

	// Example 2: Control flow
	if err := runGood(os.Args[1:]); err != nil {
		fmt.Printf("Handled argument error gracefully: %v\n", err)
	} else {
		fmt.Println("Program ran successfully!")
	}

	// This would crash the program if uncommented!
	// fmt.Println(getGreetingBad(123))
	// runBad([]string{})
}
