# Go Error Handling Guidelines

This repository serves as a reference guide for handling errors in Go idiomatically, aligning with core Go principles, the Uber Go Style Guide, and Domain-Driven Design (DDD) practices.

## 1. Idiomatic Go Guidelines

Error handling in Go isn't an afterthought; it operates under the principle that **"errors are values."** 

- **Always Check Errors**: Errors are explicit. Always check if `err != nil`. Never ignore an error unless you explicitly document why with a comment.
- **Handle or Return, Never Both**: When an error occurs, either handle it (degrade gracefully, log it, retry) **OR** return it up the call stack with added context. If you log an error and then return it, every layer will log it, creating duplicate, noisy logs.
- **Don't Panic**: `panic` and `recover` are not meant for standard control flow. Use `panic` exclusively for irrecoverable application failure (like missing required startup configuration or nil pointer dereferences). Always return `error` values under standard operation.
- **Building the Go-nion**: Add context to errors as they bubble up using `fmt.Errorf("...: %w", err)` so developers have a semantic trace of what happened (e.g., `fetching user: query failed: connection refused`).

## 2. Go Style Guidelines

The established [Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md#errors) outlines strict rules for cleanliness and predictability when writing errors.

- **Error Naming Conventions**:
  - Sentinel values (static, reusable errors) must begin with `Err` or `err` (e.g., `ErrNotFound` or `errInternal`).
  - Custom error types (structs) must end with `Error` (e.g., `PathError` or `FormatError`).
- **Choosing the Right Error Type**:
  - If the string is static and callers need to match it using `errors.Is`, use `errors.New` and export it globally.
  - If the string requires dynamic formatting but no matching is needed, use `fmt.Errorf`.
  - If the caller needs to match dynamic/complex information (using `errors.As` or `errors.AsType`), create a custom `struct` implementing the `error` interface.
- **Error Wrapping (%w vs. %v)**:
  - Add context with `%w` if the caller should be allowed to unwrap and inspect the root cause.
  - Add context with `%v` if you intentionally want to **obfuscate** the underlying error, preventing callers from relying on internal functionality or implementation details.
  - When wrapping, keep the phrasing succinct. Do not use phrases like `"failed to"` which pile up redundantly. _(Use `"process request: %w"` instead of `"failed to process request: %w"`)._

## 3. Errors in Context of Domain-Driven Design (DDD)

In a layered architecture (like DDD), error handling strategy changes depending on which boundary you cross.

### The "Split-Brain" Pattern (Domain Boundaries)
Infrastructure resources (databases, network requests, filesystems) produce raw errors that leak sensitive details (SQL paths, tokens, library versions). We call this "leaking the kitchen core."

- **The Rule of Values**: In DDD, you wouldn't return a low-level Database DTO to the presentation layer. By the same token, **you must not return raw database errors to the presentation layer.**
- Stop the infrastructure error chain at the repository boundary.

### How to Define Domain Errors
Map low-level infrastructure errors into sanitized, domain-specific errors. 

Create a custom Domain Error struct (e.g., `CustomerError` or `InvoiceIssuanceError`) that separates internal diagnostic logging from public exposure and that has semantic meaning within the domain:

```go
type DomainError struct {
    UserMsg     string // Sanitized, friendly message for public APIs ("User not found.")
    Error       error  // The raw root cause stack trace, explicitly kept for logging
    Cause       domainErrorCause // a value of type domainErrorCause (this is an example field that could be used to differentiate between different error cases within the same domain error type)
    ...         // Add any other fields you need
}

// DomainErrorCause is an enum that represents the cause of the error. This serves as an example
// of how you could use custom error types to differentiate between different error cases within the same domain error type. Different domain error types could use different enums.
type domainErrorCause int

const (
    domainErrorCauseUnknown domainErrorCause = iota
    domainErrorCauseInvalidInput
    domainErrorCauseResourceNotFound
)

// Implement the core error interface returning the publicly safe message.
func (e *DomainError) Error() string {
    return e.UserMsg
}
```

### Layer Responsibilities
1. **Repository/Infrastructure Layer**: Catches the raw (e.g., Postgres) error. Maps it into a domain error, preserving the original error inside `Error` and defining a sanitized `UserMsg`.
2. **Domain/Service Layer**: Cares about domain logic. Uses and propagates the domain error. Logic can operate on top of custom domain error types to differentiate between different error cases.
3. **Transport/Handler Layer (e.g., HTTP REST)**: 
   - Uses `errors.As` to check if the incoming error is a `DomainError`.
   - **Logs to backend**: Records `err.Error` containing the full stack trace and raw DB secrets.
   - **Responds to user**: Serializes only `err.UserMsg` and corresponding error fields into the JSON response (e.g., HTTP 500 or 404).