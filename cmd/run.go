package cmd

import (
	"fmt"
)

func runRunMode() error {
	PrintVerbose(1, "Running in --run mode")

	if installFlag && customizeFlag {
		return fmt.Errorf("cannot specify both --install and --customize")
	}
	if !installFlag && !customizeFlag {
		return fmt.Errorf("must specify one of --install or --customize")
	}

	if installFlag {
		opts.VirtMode = "install"
	}
	if customizeFlag {
		opts.VirtMode = "customize"
	}

	PrintVerbose(1, "Mode selected: %s", opts.VirtMode)

	if err := ValidateOptions(opts); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := BuildImage(opts); err != nil {
		return fmt.Errorf("build execution failed: %w", err)
	}

	return nil
}
