package cmd

import (
	"os"
)

func RequireRoot() {
	if os.Geteuid() != 0 {
		PrintError("This operation requires root. Please run kvmage with sudo.")
		os.Exit(1)
	}
}
