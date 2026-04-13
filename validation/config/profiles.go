package config

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

// ConfigFile is the ENV property for the location of the persisted JSON Profiles file.
const ConfigFile string = "PROFILES_FILE"

// Validation is a single validation.
type Validation struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Profile is a single thread-safe validation profile.
type Profile struct {
	mutex       sync.RWMutex
	name        string
	lastUpdate  time.Time
	validations []Validation
}

// profileSnapshot is a temporary struct used for marshaling to JSON.
type profileSnapshot struct {
	Name        string       `json:"name"`
	LastUpdate  time.Time    `json:"lastUpdate"`
	Validations []Validation `json:"validations"`
}

// Profiles contains a thread-safe mapping of validation Profile(s).
//
// We don't marshal this to JSON, but use a profilesSnapshot for that.
type Profiles struct {
	mutex      sync.RWMutex
	profile    map[string]*Profile
	lastUpdate time.Time
}

// profilesSnapshot is a temporary struct used for marshaling to JSON.
type profilesSnapshot struct {
	Profile    map[string]profileSnapshot `json:"profiles"`
	LastUpdate time.Time                  `json:"lastUpdate"`
}

// NewProfile is a constructor function to initialize a new Profile.
func NewProfile(name string, validations []Validation) (*Profile, error) {
	if name == "" {
		return nil, fmt.Errorf("profile name cannot be empty")
	}

	return &Profile{
		name:        name,
		lastUpdate:  time.Now(),
		validations: validations,
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
	var refreshedProfiles profilesSnapshot

	// Get the location of the persisted Profiles file
	filePath := os.Getenv(ConfigFile)
	if filePath == "" {
		return fmt.Errorf("environment variable %s is not set or empty", ConfigFile)
	}

	// Open the persisted JSON file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open '%s' file: %w", filePath, err)
	}
	defer func() {
		_ = file.Close()
	}()

	// Decode the persisted JSON file into a profilesSnapshot
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

// String returns a string representation of the Profiles instance.
func (profiles *Profiles) String() (string, error) {
	// Lock the Profiles struct for thread-safe access
	profiles.mutex.RLock()
	defer profiles.mutex.RUnlock()

	// Take a snapshot of the current profiles for serialization
	snapshot := profiles.snapshot()

	// Serialize the snapshot to JSON
	jsonData, jsonErr := json.MarshalIndent(snapshot, "", "  ")
	if jsonErr != nil {
		return "", fmt.Errorf("failed to serialize profiles to JSON: %w", jsonErr)
	}

	return string(jsonData), nil
}

// Save the Profiles to a pre-configured JSON file path location.
func (profiles *Profiles) Save() error {
	// Get the JSON file path from the environment variable
	filePath := os.Getenv(ConfigFile)
	if filePath == "" {
		return fmt.Errorf("environment variable '%s' is not set or empty", ConfigFile)
	}

	// Ensure the directory with the JSON file exists
	dirPath := filepath.Dir(filePath)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory '%s': %w", dirPath, err)
	}

	// Lock the Profiles struct for thread-safe access
	profiles.mutex.RLock()
	defer profiles.mutex.RUnlock()

	// Take a snapshot of the current profiles for serialization
	snapshot := profiles.snapshot()

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

	validations := make([]string, 0, len(profile.validations))

	for _, validation := range profile.validations {
		validations = append(validations, validation.Name)
	}

	return append([]string(nil), validations...)
}

// AddValidation adds a new validation name to the current Profile.
func (profile *Profile) AddValidation(name string, description string) {
	profile.mutex.Lock()
	defer profile.mutex.Unlock()

	for _, profileValidation := range profile.validations {
		if profileValidation.Name == name {
			return // Skip duplicates
		}
	}

	profile.validations = append(profile.validations, Validation{name, description})
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
func (profile *Profile) SetValidations(validations []Validation) {
	profile.mutex.Lock()
	defer profile.mutex.Unlock()

	profile.lastUpdate = time.Now()
	uniqueValidations := make(map[Validation]struct{})

	for _, validation := range validations {
		uniqueValidations[validation] = struct{}{}
	}

	profile.validations = make([]Validation, 0, len(uniqueValidations))
	for validation := range uniqueValidations {
		profile.validations = append(profile.validations, validation)
	}

	// Sort for consistency
	sort.Sort(ByName(profile.validations))
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

// ByName implements sort.Interface for []Validation based on the Name field.
//
// This is used to sort the validations in a Profile by name.
//
// See https://golang.org/pkg/sort/#example_Sort_intSlice for more details.
//
// This is used to sort the validations in a Profile by name.
//
// See https://golang.org/pkg/sort/#example_Sort_intSlice for more details.
type ByName []Validation

func (v ByName) Len() int           { return len(v) }
func (v ByName) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v ByName) Less(i, j int) bool { return v[i].Name < v[j].Name }

// snapshot provides a copy of Profile to marshal to JSON.
func (profile *Profile) snapshot() profileSnapshot {
	profile.mutex.RLock()
	defer profile.mutex.RUnlock()

	// Populate the temporary struct with current values
	return profileSnapshot{
		Name:        profile.name,
		LastUpdate:  profile.lastUpdate,
		Validations: append([]Validation{}, profile.validations...),
	}
}

// snapshot provides a copy of Profiles to marshal to JSON.
func (profiles *Profiles) snapshot() profilesSnapshot {
	profiles.mutex.RLock()
	defer profiles.mutex.RUnlock()

	if profiles.profile == nil {
		return profilesSnapshot{
			Profile:    make(map[string]profileSnapshot),
			LastUpdate: time.Now(),
		}
	}

	// Populate the temporary struct with current values
	snapshot := profilesSnapshot{
		Profile:    make(map[string]profileSnapshot),
		LastUpdate: profiles.lastUpdate,
	}

	for name, profile := range profiles.profile {
		snapshot.Profile[name] = profile.snapshot()
	}

	return snapshot
}
