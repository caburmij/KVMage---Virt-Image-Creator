package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kvmage",
	Short: "KVMage - A lightweight image builder for KVM",
	Long:  "KVMage builds and customizes QCOW images using virt-install, virt-customize, and related tools.",
	RunE:  runMain,
}

func init() {
	rootCmd.SetHelpFunc(CustomHelp)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		PrintError("%v", err)
		os.Exit(1)
	}
}
