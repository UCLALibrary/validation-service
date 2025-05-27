package util

import (
	"encoding/json"
	"fmt"
)

func ExampleProfiles() {
	profiles := NewProfiles()

	profile, err := NewProfile("example", []string{"Validation1", "Validation2"})
	if err != nil {
		panic(err)
	}

	if err := profiles.SetProfile(profile); err != nil {
		panic(err)
	}

	// Print the profile name
	fmt.Println(profiles.GetProfile("example").GetName())

	// Print the snapshot JSON
	snapshot := profiles.Snapshot()
	jsonData, err := json.Marshal(snapshot)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonData))
}
