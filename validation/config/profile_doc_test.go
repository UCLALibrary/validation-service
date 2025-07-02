package config

import (
	"encoding/json"
	"fmt"
)

func ExampleProfiles() {
	profiles := NewProfiles()

	profile, err := NewProfile("example", []Validation{
		{"Validation1", "Validation 1 description"},
		{"Validation2", "Validation 2 description"},
	})
	if err != nil {
		panic(err)
	}

	if err := profiles.SetProfile(profile); err != nil {
		panic(err)
	}

	// Print the profile name
	fmt.Println(profiles.GetProfile("example").GetName())

	// Print the snapshot JSON
	snapshot := profiles.snapshot()
	bytes, err := json.Marshal(snapshot)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))
}
