# Contextual Errors Experiment

## Goal
The goal of this experiment is to explore a mechanism for enriching Go errors with structured context (key-value metadata) as they bubble up the call stack, without repeatedly mutating error strings (which often leads to anti-patterns like `"failed to do X: failed to do Y: cause"`). Additionally, it specifically aims to solve the **"interface trap"** that commonly occurs when returning custom error pointers as the standard `error` interface.

## Implementation Details
1. **`contextualError` Struct**: A custom error type (`pkg/contextual_errors/contextual_error.go`) that wraps an underlying `error` and holds a `map[string]any` representing the metadata context.
2. **Builder Pattern API**: Exposes `With()` and `Wrap()` functions to initialize the contextual error. Methods like `.Str()`, `.Int()`, `.Bool()`, and `.Any()` enable a seamless, chainable API for adding metadata.
3. **Avoiding the Interface Trap**: The `.Err()` method finalizes the builder chain and explicitly returns the standard `error` interface. If the tracked error is `nil`, it returns a pure untyped `nil` `error` instead of a non-nil interface pointing to a nil `*contextualError`, preventing false positives in `if err != nil` checks.
4. **Metadata Extraction (`GetMetadata`)**: Uses `errors.Unwrap` to walk down the entire error chain, aggregating all metadata maps encountered. This single flat map can then be directly passed to structured logging solutions (like `zerolog`) at the outermost layer (e.g., an HTTP handler) for complete context visibility.
