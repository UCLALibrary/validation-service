// Package validation provides tools to validate CSV data.
package validation

import (
	"fmt"
	"os"

	"github.com/UCLALibrary/validation-service/validation/config"

	"github.com/UCLALibrary/validation-service/validation/csv"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Engine performs the CSV file validations.
type Engine struct {
	logger   *zap.Logger
	registry *Registry
	profiles *config.Profiles
}

// NewEngine creates a new validation engine.
func NewEngine(suppliedLogger ...*zap.Logger) (*Engine, error) {
	var logger *zap.Logger
	var err error

	// Use a supplied logger or create our own from scratch
	if len(suppliedLogger) != 0 {
		logger = suppliedLogger[0] // We just take the first one, and ignore any others
	} else {
		logger, err = buildLogger()
		if err != nil {
			return nil, err
		}
	}
	defer func() {
		if syncErr := logger.Sync(); syncErr != nil {
			// Combine the deferred sync error with a pre-existing existing error
			if err == nil {
				err = fmt.Errorf("error syncing logger: %w", syncErr)
			} else {
				err = fmt.Errorf("error syncing logger: %v; %w", syncErr, err)
			}
		}
	}()

	// Create a new Profiles instance and load its persisted data from disk
	profiles := config.NewProfiles()
	if err = profiles.Refresh(); err != nil {
		return nil, fmt.Errorf("failed to refresh profiles: %w", err)
	}

	// Create a new validations registry so we can retrieve CSV validations
	registry, regErr := NewRegistry(profiles, logger)
	if regErr != nil {
		return nil, regErr
	}

	// Else, return a newly constructed engine
	return &Engine{
		logger:   logger,
		registry: registry,
		profiles: profiles,
	}, nil
}

// GetLogger gets the logger used by the validation engine.
func (engine *Engine) GetLogger() *zap.Logger {
	return engine.logger
}

// GetValidators returns just the validators that are associated with the supplied profile names, or all validators
// if no profile names are passed as arguments.
func (engine *Engine) GetValidators(profileNames ...string) ([]Validator, error) {
	var checks []Validator

	// If no profiles are requested, return all the validators
	if len(profileNames) == 0 {
		validators, err := engine.registry.GetValidators(nil, engine.profiles)
		if err != nil {
			return nil, err
		}

		// We only care about the validators at this point, not their names
		return validators.Checks, nil
	}

	// Keep a record of added validator names (since profiles might have duplicates)
	existing := make(map[string]struct{})

	// For each profile we've passed in... (the most common case will just be one)
	for _, profileName := range profileNames {
		profile := engine.profiles.GetProfile(profileName)

		if profile != nil {
			validations := removeExisting(profile.GetValidations(), existing)

			// Right now, we're just passing Profiles as arguments to validator constructors
			validators, err := engine.registry.GetValidators(validations, engine.profiles, engine.logger)
			if err != nil {
				return nil, err
			}

			// In the case of profiles containing the same checks, we only want to add a check once
			for index, validatorName := range validators.Names {
				if _, exists := existing[validatorName]; !exists {
					existing[validatorName] = struct{}{}
					checks = append(checks, validators.Checks[index])
				}
			}
		}
	}

	return checks, nil
}

// Validate validates the supplied CSV data with the supplied profile name in mind.
func (engine *Engine) Validate(profile string, csvData [][]string) error {
	var errs error

	validators, err := engine.GetValidators(profile)
	if err != nil {
		return fmt.Errorf("failed to get validators: %w", err)
	}

	// Check to see if we have validators associated with the supplied profile
	if len(validators) == 0 {
		return fmt.Errorf("no validators found for profile: %s", profile)
	}

	// Have each validator check each cell in the supplied csvData
	for _, validator := range validators {
		for rowIndex, row := range csvData {
			for colIndex := range row {
				// Validate the data cell we're on, passing the entire CSV data matrix for additional context
				err := validator.Validate(profile, csv.Location{RowIndex: rowIndex, ColIndex: colIndex}, csvData)
				if err != nil {
					errs = multierr.Combine(errs, err)
				}
			}
		}
	}

	return errs
}

// removeExisting removes validations from a supplied slice if they already exist in the supplied map.
func removeExisting(validations []string, existing map[string]struct{}) []string {
	newValidations := make([]string, 0, len(validations)) // Constrain by max size

	for _, validation := range validations {
		if _, found := existing[validation]; !found {
			newValidations = append(newValidations, validation)
		}
	}

	return newValidations
}

// buildLogger constructs a logger for use (if one wasn't supplied to the engine at initialization).
func buildLogger() (*zap.Logger, error) {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Sampling = nil // We want to see all the things we log at a given level

	// Explicitly set logs to be written to stdout instead of stderr
	loggerConfig.OutputPaths = []string{"stdout"}
	loggerConfig.ErrorOutputPaths = []string{"stderr"}

	// In the production code, we just care about the ENV settings, not arg flags
	logLevel := os.Getenv("LOG_LEVEL")

	// Set the log level based on the flag or environment variable
	switch logLevel {
	case "debug":
		loggerConfig.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		loggerConfig.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		loggerConfig.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		loggerConfig.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		loggerConfig.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	// Build the logger from our specific configuration
	logger, err := loggerConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	return logger, nil
}
