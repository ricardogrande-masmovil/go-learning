---
name: go-error-style
description: Refactors Go code to comply with the project's strict error handling guidelines and Domain-Driven Design (DDD) error boundaries. Use this skill when asked to review or refactor Go error handling, or when creating new error types in a DDD architecture.
metadata:
  version: "1.0"
---

# Go Error Handling Refactoring Guide

## Purpose
This skill ensures that Go code strictly adheres to idiomatic error handling, the Uber Go Style Guide, and the project's Domain-Driven Design (DDD) "Split-Brain" error patterns.

## When to Use
Trigger this skill whenever:
- Refactoring existing Go code that returns or handles errors.
- Implementing new repository, service, or transport/HTTP handler layers in a DDD architecture.
- Reviewing code for error leaks (e.g., exposing raw database errors to APIs).
- The user requests to "fix error handling", "apply error codestyle", or "implement DDD errors".

## Core Rules to Enforce

1. **Handle or Return (Never Both)**: Ensure errors are either handled (e.g., degraded gracefully, logged, retried) OR returned up the stack. Never log and return the same error.
2. **Go-nion Context Wrapping**: Ensure contextual wrapping is succinct. Remove redundant phrases like `"failed to"`.
   - Use `%w` to expose the root cause to callers.
   - Use `%v` to intentionally obfuscate the root cause from callers.
3. **No Panics**: Replace `panic()` with returned `error` values for standard control flow.
4. **Error Naming**:
   - Variables (sentinels): Must start with `Err` or `err` (e.g., `ErrNotFound`).
   - Types (structs): Must end with `Error` (e.g., `PathError`).
5. **The Split-Brain Pattern (DDD)**:
   - **Infrastructure layer** (Repositories) must map raw/external errors (like DB timeouts, SQL errors) into Domain Error structs before returning them. Never leak raw external errors to the domain.
   - **Domain errors** must be custom structs (e.g., `CustomerError`, `InvoiceIssuanceError`) containing at least:
     - `UserMsg string`: Sanitized error for public APIs.
     - `Error error`: The raw root cause stack trace (used ONLY for internal backend logging).
     - `Cause domainErrorCause`: A custom enum to differentiate error cases for domain logic branching.
   - **Transport layer** (HTTP) must log the `err.Error` internally, but serialize and return only the `err.UserMsg` and contextual enum values to the client.
6. **Respect Existing Logging Standards**: While applying these error rules, DO NOT override the project's existing logging formats or libraries. Adapt the error handling logic to fit within their established logging paradigms (e.g., using specific structured fields in `zap`, `logrus`, `zerolog`, etc).

## Transformation Examples

### Example 1: Redundant Wrapping & Logging (One-Shot)

**Bad (Input):**
```go
func getUser(id string) error {
    err := db.Fetch(id)
    if err != nil {
        log.Printf("failed to get user: %v", err)
        return fmt.Errorf("failed to get user: %w", err)
    }
    return nil
}
```

**Good (Output):**
```go
func getUser(id string) error {
    err := db.Fetch(id)
    if err != nil {
        // Removed logging statement to prevent double-logging.
        // Removed "failed to" anti-pattern for succinct wrapping.
        return fmt.Errorf("get user: %w", err)
    }
    return nil
}
```

### Example 2: Implementing a DDD Split-Brain Error Boundary (Multi-Shot)

This is a complex transformation requiring definition of the domain error, the infrastructure mapping, and the transport handling.

**Bad (Input):**
```go
// Infra
func (r *Repo) CreateUser(u *User) error {
    return db.Exec("INSERT...", u) // Returns raw postgres error (e.g. unique constraint violation)
}

// Service
func (s *Service) Register(u *User) error {
    return s.repo.CreateUser(u) // Leaks raw error
}

// HTTP Handler
func (h *Handler) PostUser(w http.ResponseWriter, r *http.Request) {
    if err := h.svc.Register(user); err != nil {
        // Leaking DB details to the public API
        http.Error(w, err.Error(), 500) 
    }
}
```

**Good (Output):**
```go
// 1. Define the Domain Error & Causes
type userErrorCause int

const (
    userErrorCauseUnknown userErrorCause = iota
    userErrorCauseDuplicateInput
)

type UserRegistrationError struct {
    UserMsg string
    Error   error
    Cause   userErrorCause
}

func (e *UserRegistrationError) Error() string { return e.UserMsg }

// 2. Map at the Infrastructure / Repository Boundary
func (r *Repo) CreateUser(u *User) error {
    err := db.Exec("INSERT...", u)
    if err != nil {
        // Stop the chain. Map to Domain Error.
        cause := userErrorCauseUnknown
        msg := "We could not register the user at this time."
        
        if isUniqueConstraint(err) {
            cause = userErrorCauseDuplicateInput
            msg = "A user with these details already exists."
        }
        
        return &UserRegistrationError{
            UserMsg: msg,
            Error:   fmt.Errorf("db insert user %s: %w", u.ID, err),
            Cause:   cause,
        }
    }
    return nil
}

// 3. Service Layer Propagates the Domain Error
func (s *Service) Register(u *User) error {
    err := s.repo.CreateUser(u)
    // Could branch logic based on `err.Cause` here if needed by domain
    return err 
}

// 4. Transport Layer Logs the Root, Returns the Sanitized Message
func (h *Handler) PostUser(w http.ResponseWriter, r *http.Request) {
    if err := h.svc.Register(user); err != nil {
        var domainErr *UserRegistrationError
        if errors.As(err, &domainErr) {
            log.Printf("[SYSTEM ERR] User registration failed: %v", domainErr.Error)
            
            code := http.StatusInternalServerError
            if domainErr.Cause == userErrorCauseDuplicateInput {
                code = http.StatusConflict
            }
            http.Error(w, domainErr.UserMsg, code)
            return
        }
        // Fallback for unknown bugs
        log.Printf("[SYSTEM ERR] Unknown error: %v", err)
        http.Error(w, "Internal Server Error", 500)
    }
}
```

### Example 3: Panic to Error Idiom (One-Shot)

**Bad (Input):**
```go
func parseConfig(data []byte) *Config {
    if len(data) == 0 {
        panic("config data is empty")
    }
    // ...
}
```

**Good (Output):**
```go
func parseConfig(data []byte) (*Config, error) {
    if len(data) == 0 {
        // Replaced panic with a properly formatted sentinel/static error
        return nil, errors.New("config data is empty")
    }
    // ...
}
```

## 🛠 Required Steps When Applying

1. **Analyze Boundaries**: Determine if the code being refactored crosses an infrastructure/domain boundary. 
2. **Remove Anti-Patterns**: Scan for `panic()`, any logging statement (e.g., `log.Printf`, `logrus.Error`, `zap.Error`) followed by `return`, and `"failed to"` string wrapping.
3. **Apply DDD Structs**: If an infrastructure error is leaking, create the `[Domain]Error` struct with `UserMsg`, `Error`, and `Cause` fields.
4. **Update Handlers**: Ensure the HTTP/Transport layer uses `errors.As` or `errors.AsType` to trap the domain error, log `.Error` (preserving any project-specific logging standards like JSON tags or log levels), and respond with `.UserMsg`.
