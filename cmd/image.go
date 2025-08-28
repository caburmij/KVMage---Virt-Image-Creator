package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func FinalizeImage(opts *Options) error {
	RequireRoot()

	src := TempImagePath
	dstDir := opts.ImageDestination
	imageName := opts.ImageName + ".qcow2"
	dst := filepath.Join(dstDir, imageName)

	if _, err := os.Stat(dstDir); os.IsNotExist(err) {
		return fmt.Errorf("destination directory does not exist: %s", dstDir)
	}

	PrintVerbose(2, "Finalizing image using qemu-img convert: %s -> %s", src, dst)

	cmd := exec.Command("qemu-img", "convert", "-O", "qcow2", src, dst)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to finalize image using qemu-img: %w", err)
	}

	PrintVerbose(1, "Final image saved to %s", dst)
	return nil
}
