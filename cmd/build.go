package cmd

import (
	"fmt"
)

func BuildImage(opts *Options) error {
	PrintVerbose(1, "Starting image build process...")
	PrintVerbose(1, "Mode selected: %s", opts.VirtMode)

	opts.ResolvePaths()

	PrintVerbose(2, "Creating temporary image %s", opts.ImageSize)
	tempName, tempPath, err := CreateTempImage(opts)
	if err != nil {
		return fmt.Errorf("temp image creation failed: %w", err)
	}

	PrintVerbose(2, "Copying input files to temporary working directory...")
	if err := CopyInputFilesToTempDir(opts); err != nil {
		return fmt.Errorf("input file copy failed: %w", err)
	}

	switch opts.VirtMode {
	case "install":
		PrintVerbose(1, "Running install mode...")
		if err := RunInstall(opts, tempName, tempPath); err != nil {
			return fmt.Errorf("install failed: %w", err)
		}

		PrintVerbose(2, "Verifying that OS installation completed successfully...")
		if err := VerifyInstallationComplete(tempPath); err != nil {
			return fmt.Errorf("image verification failed: %w", err)
		}
		PrintVerbose(2, "Image verification passed.")

	case "customize":
		PrintVerbose(1, "Running customize mode...")
		if err := RunCustomize(opts, tempName, tempPath); err != nil {
			return fmt.Errorf("customize failed: %w", err)
		}
	default:
		return fmt.Errorf("unknown virt_mode: %s", opts.VirtMode)
	}

	PrintVerbose(2, "Finalizing image...")
	if err := FinalizeImage(opts); err != nil {
		return fmt.Errorf("finalizing image failed: %w", err)
	}

	PrintVerbose(1, "Image build process complete.")
	return nil
}
