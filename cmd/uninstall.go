package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func RunUninstall() error {
	RequireRoot()

	fmt.Print("Are you sure you want to uninstall KVMage? [y/N]: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))

	if input != "y" && input != "yes" {
		Print("Aborted. KVMage was not uninstalled.")
		return nil
	}

	path := "/usr/local/bin/kvmage"
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to uninstall: %w", err)
	}

	Print("KVMage uninstalled successfully.")
	return nil
}
