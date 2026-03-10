package main

import (
	"errors"
	"fmt"
)

// ErrNotExist is a sentinel error for our generic infrastructure operations.
var ErrNotExist = errors.New("resource does not exist")

// CustomPathError is a structured error we might want to extract during peeling.
type CustomPathError struct {
	Path    string
	Message string
}

func (e *CustomPathError) Error() string { return fmt.Sprintf("%s: %s", e.Message, e.Path) }

func main() {
	fmt.Println("=== 03: Peeling the layers and grouping the ingredients ===")

	peelingErrors()
	fmt.Println()
	groupingErrors()
}

func peelingErrors() {
	fmt.Println("--- Inspection (Peeling) ---")

	// 1. Using errors.Is to peel layers and check for sentinel values.
	rootErr := ErrNotExist
	// GOOD: Keep context succinct. Avoid "failed to process request..." which piles up.
	wrappedErr := fmt.Errorf("process request: %w", rootErr)

	// errors.Is checks if any layer of the onion matches the specific value.
	if errors.Is(wrappedErr, ErrNotExist) {
		fmt.Println("Found ErrNotExist inside the wrapped error using errors.Is!")
	}

	// 2. Extracting custom error types.
	customRoot := &CustomPathError{Path: "/dev/zero", Message: "permission denied"}
	// GOOD: Keep context succinct. Avoid "operation failed..."
	wrappedCustom := fmt.Errorf("operation: %w", customRoot)

	// Historical correct approach (Go 1.13+):
	var custom *CustomPathError
	if errors.As(wrappedCustom, &custom) {
		fmt.Printf("Extracted CustomPathError using errors.As - Path: %s\n", custom.Path)
	}

	// New Go 1.26 feature! Type-safe inspection.
	// (Documented from slides. Requires Go 1.26+)
	if pathErr, ok := errors.AsType[*CustomPathError](wrappedCustom); ok {
		fmt.Printf("Extracted with AsType... Path: %s\n", pathErr.Path)
	}
}

func groupingErrors() {
	fmt.Println("--- Joining (Grouping) ---")

	// Sometimes an operation can fail for multiple reasons simultaneously (e.g., validation).
	valErr1 := errors.New("name is required")
	valErr2 := errors.New("age must be positive")

	// errors.Join groups multiple error values into a single struct
	// (*errors.joinError) containing a slice of errors.
	err := errors.Join(valErr1, valErr2)

	fmt.Printf("Combined Error:\n%v\n\n", err)

	// We can still inspect the group! If ANY of the joined errors match, it returns true.
	if errors.Is(err, valErr1) {
		fmt.Println("The joined error array contained the 'name is required' error.")
	}
}
