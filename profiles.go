// This code provides utilities for creating and working with validation profiles.
//
// Example usage:
//
//	profiles := NewProfiles()
//	if profile, err := NewProfile("example", []string{"Validation1", "Validation2"}); err != nil {
//		require.NoError(t, err)
//	} else {
//		err = profiles.SetProfile(profile)
//		require.NoError(t, err)
//
//		fmt.Println(profiles.GetProfile("example").GetName())
//	}
//
//	snapshot := profiles.Snapshot()
//	if jsonData, err := json.Marshal(snapshot); err != nil {
//		require.NoError(t, err)
//	} else {
//		fmt.Println(string(jsonData))
//	}
package main

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// Profile is a single thread-safe validation profile.
type Profile struct {
	mutex       sync.RWMutex
	name        string
	lastUpdate  time.Time
	validations []string
}

// ProfileSnapshot is a temporary struct used for marshaling to JSON.
type ProfileSnapshot struct {
	Name        string    `json:"name"`
	LastUpdate  time.Time `json:"lastUpdate"`
	Validations []string  `json:"validations"`
}

// Profiles contains a thread-safe mapping of validation Profile(s).
//
// We don't marshal this to JSON, but use a ProfilesSnapshot for that.
type Profiles struct {
	mutex      sync.RWMutex
	profile    map[string]*Profile
	lastUpdate time.Time
}

// ProfilesSnapshot is a temporary struct used for marshaling to JSON.
type ProfilesSnapshot struct {
	Profile    map[string]ProfileSnapshot `json:"profiles"`
	LastUpdate time.Time                  `json:"lastUpdate"`
}

// NewProfile is a constructor function to initialize a new Profile.
func NewProfile(name string, validations []string) (*Profile, error) {
	if name == "" {
		return nil, fmt.Errorf("profile name cannot be empty")
	}

	return &Profile{
		name:        name,
		lastUpdate:  time.Now(),
		validations: append([]string(nil), validations...),
	}, nil
}

// NewProfiles is a constructor function to initialize a new Profiles.
func NewProfiles() *Profiles {
	return &Profiles{
		profile: make(map[string]*Profile),
	}
}

// Count the number of Profile(s) in this Profiles instance.
func (profiles *Profiles) Count() int {
	profiles.mutex.RLock()
	defer profiles.mutex.RUnlock()
	return len(profiles.profile)
}

// GetName gets the name of the current Profile.
func (profile *Profile) GetName() string {
	profile.mutex.RLock()
	defer profile.mutex.RUnlock()
	return profile.name
}

// GetLastUpdate gets the last update of the current Profile.
func (profile *Profile) GetLastUpdate() time.Time {
	profile.mutex.RLock()
	defer profile.mutex.RUnlock()
	return profile.lastUpdate
}

// GetValidations gets the validation names of the current Profile.
func (profile *Profile) GetValidations() []string {
	profile.mutex.RLock()
	defer profile.mutex.RUnlock()
	return append([]string(nil), profile.validations...)
}

// AddValidation adds a new validation name to the current Profile.
func (profile *Profile) AddValidation(validation string) {
	profile.mutex.Lock()
	defer profile.mutex.Unlock()

	for _, v := range profile.validations {
		if v == validation {
			return // Skip duplicates
		}
	}

	profile.validations = append(profile.validations, validation)
}

// GetProfile gets the Profile with the supplied name.
func (profiles *Profiles) GetProfile(name string) *Profile {
	profiles.mutex.RLock()
	defer profiles.mutex.RUnlock()

	profile, exists := profiles.profile[name]
	if !exists {
		return nil
	}

	return profile
}

// SetName a new Profile name.
func (profile *Profile) SetName(name string) {
	profile.mutex.Lock()
	defer profile.mutex.Unlock()

	profile.lastUpdate = time.Now()
	profile.name = name
}

// SetValidations a Profile validations.
func (profile *Profile) SetValidations(validations []string) {
	profile.mutex.Lock()
	defer profile.mutex.Unlock()

	profile.lastUpdate = time.Now()
	uniqueValidations := make(map[string]struct{})
	for _, v := range validations {
		uniqueValidations[v] = struct{}{}
	}

	profile.validations = make([]string, 0, len(uniqueValidations))
	for v := range uniqueValidations {
		profile.validations = append(profile.validations, v)
	}

	// Sort for consistency
	sort.Strings(profile.validations)
}

// SetProfile sets a new Profile in Profiles.
func (profiles *Profiles) SetProfile(profile *Profile) error {
	if profile == nil {
		return fmt.Errorf("cannot set a nil profile")
	}

	if profile.GetName() == "" {
		return fmt.Errorf("cannot set a profile with an empty name")
	}

	profiles.mutex.Lock()
	defer profiles.mutex.Unlock()

	if profiles.profile == nil {
		profiles.profile = make(map[string]*Profile)
	}

	profiles.lastUpdate = time.Now()
	profiles.profile[profile.GetName()] = profile

	return nil
}

// Snapshot provides a copy of Profile to marshal to JSON.
func (profile *Profile) Snapshot() ProfileSnapshot {
	profile.mutex.RLock()
	defer profile.mutex.RUnlock()

	// Populate the temporary struct with current values
	return ProfileSnapshot{
		Name:        profile.name,
		LastUpdate:  profile.lastUpdate,
		Validations: append([]string(nil), profile.validations...),
	}
}

// Snapshot provides a copy of Profiles to marshal to JSON.
func (profiles *Profiles) Snapshot() ProfilesSnapshot {
	profiles.mutex.RLock()
	defer profiles.mutex.RUnlock()

	if profiles.profile == nil {
		return ProfilesSnapshot{
			Profile:    make(map[string]ProfileSnapshot),
			LastUpdate: time.Now(),
		}
	}

	// Populate the temporary struct with current values
	snapshot := ProfilesSnapshot{
		Profile:    make(map[string]ProfileSnapshot),
		LastUpdate: profiles.lastUpdate,
	}

	for name, profile := range profiles.profile {
		snapshot.Profile[name] = profile.Snapshot()
	}

	return snapshot
}
