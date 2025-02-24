package validation

import (
	"fmt"

	"github.com/UCLALibrary/validation-service/validation/checks"
	"github.com/UCLALibrary/validation-service/validation/config"
	"go.uber.org/zap"
)

// Registry keeps a collection of registered validators.
type Registry struct {
	profiles *config.Profiles
	logger   *zap.Logger
}

// Validators is a (sub)set of validation checks and their names.
//
// We use a struct instead of a map to preserve validator order.
type Validators struct {
	Names  []string
	Checks []Validator
}

// The constructor creates a new instance of the Validator interface.
type constructor func(args ...interface{}) (Validator, error)

// A registry of Validator constructors.
//
// New validators need to be manually configured in this map to avoid reflection.
var constructors = map[string]constructor{
	"EOLCheck": func(args ...interface{}) (Validator, error) {
		if len(args) > 0 {
			// Check if the first argument is of the type *Profiles
			if profiles, ok := args[0].(*config.Profiles); ok {
				return (&checks.EOLCheck{}).NewEOLCheck(profiles)
			}

			// EOLCheck expects *Profiles to be passed to it
			return nil, fmt.Errorf("invalid argument: expected *Profiles, found: %T", args[0])
		}

		// Default instance if no arguments are passed
		defaultProfiles := config.NewProfiles() // Assume a default constructor exists
		return (&checks.EOLCheck{}).NewEOLCheck(defaultProfiles)
	},
}

// NewRegistry creates a new registry of validators
func NewRegistry(profiles *config.Profiles, logger *zap.Logger) (*Registry, error) {
	if profiles == nil {
		return nil, fmt.Errorf("supplied Profiles cannot be nil")
	}

	if logger == nil {
		return nil, fmt.Errorf("supplied Logger cannot be nil")
	}

	return &Registry{
		profiles: profiles,
		logger:   logger,
	}, nil
}

// GetValidators gets new instances of the validators for use in the validation engine.
func (registry *Registry) GetValidators(validatorNames []string, args ...interface{}) (*Validators, error) {
	var nameCount = len(validatorNames)
	var requested = make(map[string]struct{}, nameCount)
	var validators = Validators{
		Names:  []string{},
		Checks: []Validator{},
	}

	// Put all the requested validator names in a map for easier access
	for _, name := range validatorNames {
		requested[name] = struct{}{}
	}

	// Loop through registry and return just the requested validators
	for name, constructor := range constructors {
		// len(validatorNames) will return 0 if the slice is empty or nil
		if _, exists := requested[name]; exists || nameCount == 0 {
			// We call the constructors with the args requested at GetValidators
			validator, err := constructor(args...)
			if err != nil {
				registry.logger.Error("error creating validator", zap.Error(err))
				continue // Skip this invalid validator and move to the next
			}

			// Add error-free validators to the slice we're returning
			validators.Names = append(validators.Names, name)
			validators.Checks = append(validators.Checks, validator)
		}
	}

	return &validators, nil
}
