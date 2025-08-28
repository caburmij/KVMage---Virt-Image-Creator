package cmd

var Version string = "unknown"

func PrintVersion() {
	Print("KVMage version: %s", Version)
}
