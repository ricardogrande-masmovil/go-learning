package main

import (
	"context"

	"github.com/ricardogrande-masmovil/go-learning/errors/contextual-error/internal/app"
	"github.com/ricardogrande-masmovil/go-learning/errors/contextual-error/internal/infrastructure/postgres"
	"github.com/ricardogrande-masmovil/go-learning/errors/contextual-error/internal/ports"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Info().Msg("Initializing contextual-error demonstration...")

	// 1. Initialize dependencies
	// For this demonstration, we use a nil sql.DB since the fake client doesn't use it
	repo := postgres.NewClient(nil)
	service := app.NewUserService(repo)
	handler := ports.NewHTTPHandler(service)

	ctx := context.Background()

	// 2. Demonstration: Validation Error
	log.Info().Msg("--- Scenario 1: Validation Error (Empty Name) ---")
	handler.CreateUser(ctx, "user-123", "")
	// Expected: HTTP 500 log containing validation_cause: "invalid_name"

	// 3. Demonstration: Infrastructure Error (Primary Key Violation)
	log.Info().Msg("--- Scenario 2: Infrastructure Error on Create ---")
	handler.CreateUser(ctx, "user-123", "Alice")
	// Expected: HTTP 500 log containing user_id "user-123" and cause "primary key violation"

	// 4. Demonstration: Infrastructure Error (Connection Refused)
	log.Info().Msg("--- Scenario 3: Infrastructure Error on Get ---")
	handler.GetUserByID(ctx, "user-123")
	// Expected: HTTP 500 log containing user_id "user-123" and cause "connection refused"

	log.Info().Msg("Demonstration finished.")
}
