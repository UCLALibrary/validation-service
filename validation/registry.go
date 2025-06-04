package validation

import (
	"fmt"

	"github.com/UCLALibrary/validation-service/validation/util"

	"github.com/UCLALibrary/validation-service/validation/checks"
	"go.uber.org/zap"
)

// Registry keeps a collection of registered validators.
type Registry struct {
	profiles *util.Profiles
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
			if profiles, ok := args[0].(*util.Profiles); ok {
				return checks.NewEOLCheck(profiles)
			}

			// EOLCheck expects *Profiles to be passed to it
			return nil, fmt.Errorf("invalid argument: expected *Profiles, found: %T", args[0])
		}

		// Default instance if no arguments are passed
		defaultProfiles := util.NewProfiles() // Assume a default constructor exists
		return checks.NewEOLCheck(defaultProfiles)
	},
	"ARKCheck": func(args ...interface{}) (Validator, error) {
		if len(args) > 0 {
			// Check if the first argument is of the type *Profiles
			if profiles, ok := args[0].(*util.Profiles); ok {
				return checks.NewARKCheck(profiles)
			}

			// ARKCheck expects *Profiles to be passed to it
			return nil, fmt.Errorf("invalid argument: expected *Profiles, found: %T", args[0])
		}

		// Default instance if no arguments are passed
		defaultProfiles := util.NewProfiles() // Assume a default constructor exists
		return checks.NewEOLCheck(defaultProfiles)
	},
	"LicenseCheck": func(args ...interface{}) (Validator, error) {
		if len(args) > 0 {
			// Check if the first argument is of the type *Profiles
			if profiles, ok := args[0].(*util.Profiles); ok {
				return checks.NewLicenseCheck(profiles)
			}

			// LicenseCheck expects *Profiles to be passed to it
			return nil, fmt.Errorf("invalid argument: expected *Profiles, found: %T", args[0])
		}

		// Default instance if no arguments are passed
		defaultProfiles := util.NewProfiles() // Assume a default constructor exists
		return checks.NewLicenseCheck(defaultProfiles)
	},
	"ReqFieldCheck": func(args ...interface{}) (Validator, error) {
		if len(args) > 0 {
			var profiles *util.Profiles
			var logger *zap.Logger
			var ok bool

			// Check if the first argument is of the type *Profiles
			if profiles, ok = args[0].(*util.Profiles); !ok {
				return nil, fmt.Errorf("invalid argument: expected *util.Profiles, found: %T", args[0])
			}

			if logger, ok = args[1].(*zap.Logger); !ok {
				return nil, fmt.Errorf("invalid argument: expected *zap.Logger, found: %T", args[1])
			}

			return checks.NewReqFieldCheck(profiles, logger)
		}

		return nil, fmt.Errorf("invalid argument: expected *util.Profiles and *zap.Logger; neither were found")
	},
	"FilePathCheck": func(args ...interface{}) (Validator, error) {
		if len(args) > 0 {
			// Check if the first argument is of the type *Profiles
			if profiles, ok := args[0].(*util.Profiles); ok {
				return checks.NewFilePathCheck(profiles)
			}

			// FilePathCheck expects *Profiles to be passed to it
			return nil, fmt.Errorf("invalid argument: expected *Profiles, found: %T", args[0])
		}

		// Default instance if no arguments are passed
		defaultProfiles := util.NewProfiles()
		return checks.NewFilePathCheck(defaultProfiles)
	},
	"ObjectTypeCheck": func(args ...interface{}) (Validator, error) {
		if len(args) > 0 {
			// Check if the first argument is of the type *Profiles
			if profiles, ok := args[0].(*util.Profiles); ok {
				return checks.NewObjTypeCheck(profiles)
			}

			// ObjectTypeCheck expects *Profiles to be passed to it
			return nil, fmt.Errorf("invalid argument: expected *Profiles, found: %T", args[0])
		}

		// Default instance if no arguments are passed
		defaultProfiles := util.NewProfiles()
		return checks.NewObjTypeCheck(defaultProfiles)
	},
	"ItemSeqCheck": func(args ...interface{}) (Validator, error) {
		if len(args) > 0 {
			// Check if the first argument is of the type *Profiles
			if profiles, ok := args[0].(*util.Profiles); ok {
				return checks.NewItemSeqCheck(profiles)
			}

			// ItemSeqCheck expects *Profiles to be passed to it
			return nil, fmt.Errorf("invalid argument: expected *Profiles, found: %T", args[0])
		}

		// Default instance if no arguments are passed
		defaultProfiles := util.NewProfiles()
		return checks.NewItemSeqCheck(defaultProfiles)
	},
	"VisibilityCheck": func(args ...interface{}) (Validator, error) {
		if len(args) > 0 {
			// Check if the first argument is of the type *Profiles
			if profiles, ok := args[0].(*util.Profiles); ok {
				return checks.NewVisibilityCheck(profiles)
			}

			// VisibilityCheck expects *Profiles to be passed to it
			return nil, fmt.Errorf("invalid argument: expected *Profiles, found: %T", args[0])
		}

		// Default instance if no arguments are passed
		defaultProfiles := util.NewProfiles()
		return checks.NewVisibilityCheck(defaultProfiles)
	},
	"UnicodeCheck": func(args ...interface{}) (Validator, error) {
		if len(args) > 0 {
			// Check if the first argument is of the type *Profiles
			if profiles, ok := args[0].(*util.Profiles); ok {
				return checks.NewUnicodeCheck(profiles)
			}

			// UnicodeCheck expects *Profiles to be passed to it
			return nil, fmt.Errorf("invalid argument: expected *Profiles, found: %T", args[0])
		}

		// Default instance if no arguments are passed
		defaultProfiles := util.NewProfiles()
		return checks.NewUnicodeCheck(defaultProfiles)
	},
	"FileNameCheck": func(args ...interface{}) (Validator, error) {
		if len(args) > 0 {
			// Check if the first argument is of the type *Profiles
			if profiles, ok := args[0].(*util.Profiles); ok {
				return checks.NewFileNameCheck(profiles)
			}

			// UnicodeCheck expects *Profiles to be passed to it
			return nil, fmt.Errorf("invalid argument: expected *Profiles, found: %T", args[0])
		}

		// Default instance if no arguments are passed
		defaultProfiles := util.NewProfiles()
		return checks.NewFileNameCheck(defaultProfiles)
	},
	"MediaMetaCheck": func(args ...interface{}) (Validator, error) {
		if len(args) > 0 {
			// Check if the first argument is of the type *Profiles
			if profiles, ok := args[0].(*util.Profiles); ok {
				return checks.NewMediaMetaCheck(profiles)
			}

			// MediaMetaCheck expects *Profiles to be passed to it
			return nil, fmt.Errorf("invalid argument: expected *Profiles, found: %T", args[0])
		}

		// Default instance if no arguments are passed
		defaultProfiles := util.NewProfiles()
		return checks.NewMediaMetaCheck(defaultProfiles)
	},
}

// NewRegistry creates a new registry of validators
func NewRegistry(profiles *util.Profiles, logger *zap.Logger) (*Registry, error) {
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
