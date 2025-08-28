package cmd

import (
	"os"
	"os/signal"
	"syscall"
)

func SetupSignalHandler(cleanupFunc func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		PrintError("\nInterrupt received. Cleaning up...")
		cleanupFunc()
		os.Exit(1)
	}()
}
