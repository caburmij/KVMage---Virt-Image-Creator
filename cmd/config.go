package cmd

import (
	"fmt"
)

func runConfigMode(configPath string) error {
	profiles, err := LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	for name, opts := range profiles {
		opts.Profile = name

		PrintVerbose(1, "Building profile: %s", name)

		if err := ValidateOptions(opts); err != nil {
			return fmt.Errorf("validation failed for profile '%s': %w", name, err)
		}
		if err := BuildImage(opts); err != nil {
			return fmt.Errorf("build failed for profile '%s': %w", name, err)
		}
	}

	return nil
}
