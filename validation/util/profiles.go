// Package util provides useful resources and utilities.
//
// This file contains utilities for creating and working with validation profiles.
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
package util

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// LogLevel is the ENV property for the validation engine's configurable logging level.
const LogLevel string = "LOG_LEVEL"

// ProfilesFile is the ENV property for the location of the persisted JSON Profiles file.
const ProfilesFile string = "PROFILES_FILE"

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

// Refresh Profiles from the last persisted version on disk.
//
// This overwrites the current in-memory values.
func (profiles *Profiles) Refresh() error {
	var refreshedProfiles ProfilesSnapshot

	// Get the location of the persisted Profiles file
	filePath := os.Getenv(ProfilesFile)
	if filePath == "" {
		return fmt.Errorf("environment variable %s is not set or empty", ProfilesFile)
	}

	// Open the persisted JSON file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open '%s' file: %w", filePath, err)
	}
	//noinspection GoUnhandledErrorResult
	defer file.Close()

	// Decode the persisted JSON file into a ProfilesSnapshot
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&refreshedProfiles); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	// Check that our JSON file actually has some profiles
	if len(refreshedProfiles.Profile) == 0 {
		return fmt.Errorf("no profiles found in refreshed data")
	}

	// Create a new temporary map for refreshed *Profile(s)
	tempMap := make(map[string]*Profile)

	// Convert ProfileSnapshot(s) to *Profile(s)
	for _, refreshedProfile := range refreshedProfiles.Profile {
		profile, err := NewProfile(refreshedProfile.Name, refreshedProfile.Validations)
		if err != nil {
			return fmt.Errorf("failed to create new profile '%s': %w", refreshedProfile.Name, err)
		}

		// Check to see if our tempMap already has a Profile with the same name
		profileName := profile.GetName()
		if _, exists := tempMap[profileName]; exists {
			return fmt.Errorf("profile '%s' already exists", profileName)
		}

		tempMap[profileName] = profile
	}

	profiles.mutex.Lock()
	defer profiles.mutex.Unlock()

	// Update the Profiles map and set a new lastUpdate time
	profiles.profile = tempMap
	profiles.lastUpdate = refreshedProfiles.LastUpdate

	return nil
}

// Save the Profiles to a pre-configured JSON file path location.
func (profiles *Profiles) Save() error {
	// Get the JSON file path from the environment variable
	filePath := os.Getenv(ProfilesFile)
	if filePath == "" {
		return fmt.Errorf("environment variable '%s' is not set or empty", ProfilesFile)
	}

	// Ensure the directory with the JSON file exists
	dirPath := filepath.Dir(filePath)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory '%s': %w", dirPath, err)
	}

	// Lock the profiles struct for thread-safe access
	profiles.mutex.RLock()
	defer profiles.mutex.RUnlock()

	// Take a snapshot of the current profiles for serialization
	snapshot := profiles.Snapshot()

	// Serialize the snapshot to JSON
	jsonData, jsonErr := json.MarshalIndent(snapshot, "", "  ")
	if jsonErr != nil {
		return fmt.Errorf("failed to serialize profiles to JSON: %w", jsonErr)
	}

	// Write to a temporary file first for atomic updates
	tempFile, tempFileErr := os.CreateTemp("", "profile-*.json")
	if tempFileErr != nil {
		return fmt.Errorf("failed to create temporary file '%s': %w", tempFile.Name(), tempFileErr)
	}

	// Write the JSON data to the temporary file
	if _, err := tempFile.Write(jsonData); err != nil {
		return fmt.Errorf("failed to write JSON data to temporary file '%s': %w", tempFile.Name(), err)
	}

	// Close the temporary file
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("failed to close temporary file '%s': %w", tempFile.Name(), err)
	}

	// Rename the temporary file to the final file
	if err := os.Rename(tempFile.Name(), filePath); err != nil {
		return fmt.Errorf("failed to rename temporary file '%s' to '%s': %w", tempFile.Name(), filePath, err)
	}

	return nil
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

	for _, profileValidation := range profile.validations {
		if profileValidation == validation {
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

	for _, validation := range validations {
		uniqueValidations[validation] = struct{}{}
	}

	profile.validations = make([]string, 0, len(uniqueValidations))
	for validation := range uniqueValidations {
		profile.validations = append(profile.validations, validation)
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
