package main

import "kvmage/cmd"

func main() {
	defer cmd.CleanupArtifacts()
	cmd.Execute()
}
