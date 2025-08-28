package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func RunCustomize(opts *Options, tempName, tempPath string) error {
	RequireRoot()

	PrintVerbose(1, "Starting customize mode...")
	PrintVerbose(2, "Customizing image at path: %s", tempPath)

	// Optional: Resize the image if image_size is specified
	if opts.ImageSize != "" {
		PrintVerbose(1, "Checking image size...")

		cmd := exec.Command("qemu-img", "info", "--output=json", tempPath)
		output, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("failed to get image info: %w", err)
		}

		var info struct {
			VirtualSize int64 `json:"virtual-size"`
		}
		if err := json.Unmarshal(output, &info); err != nil {
			return fmt.Errorf("failed to parse qemu-img info output: %w", err)
		}
		currentSizeGB := float64(info.VirtualSize) / (1024 * 1024 * 1024)

		targetSizeStr := strings.TrimSuffix(opts.ImageSize, "G")
		targetSizeGB, err := strconv.ParseFloat(targetSizeStr, 64)
		if err != nil {
			return fmt.Errorf("invalid image_size format: %w", err)
		}

		if targetSizeGB < currentSizeGB {
			PrintError("Refusing to shrink image: requested size %.1f GB is less than current size %.1f GB", targetSizeGB, currentSizeGB)
			return fmt.Errorf("image_size must not be smaller than the current image size")
		} else if targetSizeGB > currentSizeGB {
			PrintVerbose(1, "Expanding image from %.1f GB to %.1f GB", currentSizeGB, targetSizeGB)
			if err := exec.Command("qemu-img", "resize", tempPath, opts.ImageSize).Run(); err != nil {
				return fmt.Errorf("qemu-img resize failed: %w", err)
			}
		} else {
			PrintVerbose(2, "Image is already at requested size: %.1f GB", currentSizeGB)
		}
	}

	// Optional: Expand partition if --image-part is set
	if opts.ImagePartition != "" {
		PrintVerbose(1, "Expanding partition: %s", opts.ImagePartition)

		backupPath := tempPath + ".orig"
		if err := os.Rename(tempPath, backupPath); err != nil {
			return fmt.Errorf("failed to rename original image for virt-resize: %w", err)
		}

		args := []string{
			"--expand", opts.ImagePartition,
			backupPath,
			tempPath,
		}

		PrintVerbose(2, "Running virt-resize with args: virt-resize %s", joinArgs(args))
		if err := exec.Command("virt-resize", args...).Run(); err != nil {
			return fmt.Errorf("virt-resize failed: %w", err)
		}
	}

	// Run virt-customize
	args := []string{"-a", tempPath}
	if verboseLevel >= 1 {
		args = append(args, "-v")
	}
	if verboseLevel >= 2 {
		args = append(args, "-x")
	}
	if TempCustomScript != "" {
		args = append(args, "--run", TempCustomScript)
		PrintVerbose(2, "Including custom script: %s", TempCustomScript)
	}
	if opts.Hostname != "" {
		args = append(args, "--hostname", opts.Hostname)
		PrintVerbose(2, "Setting hostname: %s", opts.Hostname)
	}

	PrintVerbose(3, "Running virt-customize with args: virt-customize %s", joinArgs(args))
	PrintVerbose(1, "Executing virt-customize...")

	run := func(cmdName string, args []string) error {
		cmd := exec.Command(cmdName, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	if err := run("virt-customize", args); err != nil {
		return fmt.Errorf("virt-customize failed: %w", err)
	}

	return nil
}
