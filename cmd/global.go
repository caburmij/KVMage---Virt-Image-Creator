package cmd

import (
	"fmt"
	"os"
)

func HandleGlobalFlags() bool {
	if showVersion {
		PrintVersion()
		os.Exit(0)
	}

	if uninstall {
		if err := RunUninstall(); err != nil {
			fmt.Fprintf(os.Stderr, "Uninstall failed: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if cleanupOnly {
		CleanupOrphanedTempFiles()
		os.Exit(0)
	}

	return false
}
