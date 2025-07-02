//go:build unit

package validation

import (
	"errors"
	"github.com/UCLALibrary/validation-service/validation/config"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

// TestNewRegistry creates a new registry for validators.
func TestNewRegistry(t *testing.T) {
	logger := zaptest.NewLogger(t)
	profiles := config.NewProfiles()

	// Should successfully create a registry
	reg, err := NewRegistry(profiles, logger)
	assert.NoError(t, err)
	assert.NotNil(t, reg)

	// Should return an error when profiles is nil
	reg, err = NewRegistry(nil, logger)
	assert.Error(t, err)
	assert.Nil(t, reg)

	// Should return an error when logger is nil
	reg, err = NewRegistry(profiles, nil)
	assert.Error(t, err)
	assert.Nil(t, reg)
}

// TestGetValidators tests getting the validators from the registry.
func TestGetValidators(t *testing.T) {
	logger := zaptest.NewLogger(t)
	profiles := config.NewProfiles()
	reg, err := NewRegistry(profiles, logger)
	if err != nil {
		t.Errorf("NewRegistry() error = %v", err)
	}

	// Delete map's entries so we have a fresh start to test with
	for key := range constructors {
		delete(constructors, key)
	}

	// Override constructors for controlled testing
	constructors["MockValidator"] = mockConstructor
	constructors["FailingValidator"] = failingConstructor

	t.Run("Returns all validators when validatorNames is empty", func(t *testing.T) {
		validators, err := reg.GetValidators(nil) // Request all validators
		assert.NoError(t, err)
		assert.Len(t, validators.Checks, 1) // MockValidator only, FailingValidator isn't a valid check
		assert.Len(t, constructors, 2)      // MockValidator + FailingValidator
	})

	t.Run("Returns only requested validators", func(t *testing.T) {
		validators, err := reg.GetValidators([]string{"MockValidator"})
		assert.NoError(t, err)
		assert.Len(t, validators.Checks, 1)
		assert.Equal(t, "MockValidator", validators.Names[0])
	})

	t.Run("Returns empty result if requested validator does not exist", func(t *testing.T) {
		validators, err := reg.GetValidators([]string{"NonExistentValidator"})
		assert.NoError(t, err)
		assert.Empty(t, validators.Checks)
	})

	t.Run("Handles constructor errors gracefully", func(t *testing.T) {
		validators, err := reg.GetValidators([]string{"FailingValidator"})
		assert.NoError(t, err)
		assert.Empty(t, validators.Checks) // Should skip the failing validator
	})
}

// Mock Constructor for testing that will succeed.
func mockConstructor(args ...interface{}) (Validator, error) {
	return &MockValidator{}, nil
}

// Mock constructor for testing that will fail.
func failingConstructor(args ...interface{}) (Validator, error) {
	return nil, errors.New("mock constructor error")
}
