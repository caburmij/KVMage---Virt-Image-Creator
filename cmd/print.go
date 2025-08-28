package cmd

import (
	"fmt"
	"os"
)

func Print(msg string, args ...any) {
	if quiet {
		return
	}
	fmt.Printf(msg+"\n", args...)
}

func PrintVerbose(level int, msg string, args ...any) {
	if quiet || verboseLevel < level {
		return
	}
	fmt.Printf(msg+"\n", args...)
}

func PrintError(msg string, args ...any) {
	if quiet {
		return
	}
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
}
