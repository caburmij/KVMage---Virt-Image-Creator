package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func CleanupArtifacts() {
	paths := []string{
		TempImagePath,
		TempKickstartPath,
		TempInstallMedia,
		TempImageSource,
		TempCustomScript,
	}

	for _, path := range paths {
		if path == "" {
			continue
		}

		if !isSafeTempPath(path) {
			PrintVerbose(2, "Skipping unsafe path: %s", path)
			continue
		}

		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		if err := os.Remove(path); err != nil {
			PrintError("Warning: failed to remove temp file %s: %v", path, err)
		} else {
			PrintVerbose(2, "Removed temporary file: %s", path)
		}
	}

	if TempImageName != "" {
		removeTempVM(TempImageName)
	}
}

func isSafeTempPath(path string) bool {
	return strings.HasPrefix(path, "/var/lib/libvirt/images/kvmage-")
}

func removeTempVM(vmName string) {
	PrintVerbose(2, "Checking for temporary VM: %s", vmName)

	checkCmd := exec.Command("virsh", "dominfo", vmName)
	if err := checkCmd.Run(); err != nil {
		PrintVerbose(2, "Domain %s does not exist. Skipping undefine.", vmName)
		return
	}

	PrintVerbose(2, "Undefining temporary VM: %s", vmName)

	cmd := exec.Command("virsh", "undefine", "--nvram", vmName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		PrintError("Warning: failed to undefine VM %s: %v", vmName, err)
	}
}

func CleanupOrphanedTempFiles() {
	dir := "/var/lib/libvirt/images"
	entries, err := os.ReadDir(dir)
	if err != nil {
		PrintError("Failed to read directory %s: %v", dir, err)
		os.Exit(1)
	}

	var toDelete []string
	var totalSize int64

	for _, entry := range entries {
		name := entry.Name()
		fullPath := filepath.Join(dir, name)

		if !strings.HasPrefix(name, "kvmage-") {
			continue
		}

		info, err := entry.Info()
		if err != nil || info.IsDir() {
			continue
		}

		if isFileInUse(fullPath) {
			PrintVerbose(2, "Skipping in-use file: %s", fullPath)
			continue
		}

		toDelete = append(toDelete, fullPath)
		totalSize += info.Size()
	}

	if len(toDelete) == 0 {
		Print("No orphaned kvmage temp files found.")
		return
	}

	Print("Found %d orphaned kvmage temp file(s):\n", len(toDelete))
	for _, path := range toDelete {
		size := getFileSize(path)
		Print("  %-65s %8s", path, formatSize(size))
	}
	Print("\nTotal reclaimable space: %s", formatSize(totalSize))

	Print("\nDo you want to delete these files? [y/N]: ")
	var input string
	fmt.Scanln(&input)

	if strings.ToLower(input) != "y" {
		Print("Aborted. No files deleted.")
		return
	}

	for _, path := range toDelete {
		if err := os.Remove(path); err != nil {
			PrintError("Failed to delete %s: %v", path, err)
		} else {
			Print("Deleted: %s", path)
		}
	}
}

func isFileInUse(path string) bool {
	cmd := exec.Command("lsof", "-Fn", "--", path)
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return len(output) > 0
}

func getFileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}

func formatSize(size int64) string {
	const (
		KB = 1 << 10
		MB = 1 << 20
		GB = 1 << 30
	)
	switch {
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/KB)
	default:
		return fmt.Sprintf("%d B", size)
	}
}
