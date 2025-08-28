package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func RunInstall(opts *Options, tempName, tempPath string) error {
	RequireRoot()

	PrintVerbose(1, "Starting install mode...")
	PrintVerbose(1, "Creating VM: %s", tempName)

	args := []string{
		"--name", tempName,
		"--memory", "4096",
		"--vcpus", "2",
		"--cpu", "host-passthrough",
		"--os-variant", opts.OSVariant,
		"--disk", fmt.Sprintf("path=%s,format=qcow2,bus=virtio", tempPath),
		"--location", TempInstallMedia,
		"--initrd-inject", TempKickstartPath,
		"--noreboot",
		"--wait", "-1",
	}

	PrintVerbose(2, "OS Variant: %s", opts.OSVariant)
	PrintVerbose(2, "Disk Path: %s", tempPath)
	PrintVerbose(2, "Kickstart: %s", TempKickstartPath)
	PrintVerbose(2, "Install Media: %s", TempInstallMedia)

	if opts.Network != "" {
		args = append(args, "--network", fmt.Sprintf("network=%s,model=virtio", opts.Network))
		PrintVerbose(2, "Using custom network: %s", opts.Network)
	} else {
		args = append(args, "--network", "network=default,model=virtio")
		PrintVerbose(2, "Using default network")
	}

	if opts.Firmware == "efi" {
		args = append(args, "--machine", "q35", "--boot", "uefi")
		PrintVerbose(2, "Firmware: efi (UEFI boot with Q35 chipset)")
	} else {
		PrintVerbose(2, "Firmware: bios (default)")
	}

	extraArgs := fmt.Sprintf("inst.ks=file:/%s", filepath.Base(TempKickstartPath))

	switch opts.Console {
	case "serial":
		args = append(args,
			"--graphics", "none",
			"--console", "pty,target_type=serial",
		)
		extraArgs += " console=ttyS0 inst.text"
		PrintVerbose(2, "Console type: serial")
	case "graphical":
		args = append(args,
			"--graphics", "vnc,listen=127.0.0.1",
			"--noautoconsole",
		)
		PrintVerbose(2, "Console type: graphical")
	default:
		PrintVerbose(2, "Console type: (none specified)")
	}

	args = append(args, "--extra-args", extraArgs)

	PrintVerbose(3, "Running virt-install with args: virt-install %s", joinArgs(args))

	run := func(cmdName string, args []string) error {
		cmd := exec.Command(cmdName, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		return cmd.Run()
	}

	PrintVerbose(1, "Executing virt-install for: %s", tempName)

	if err := run("virt-install", args); err != nil {
		return fmt.Errorf("virt-install failed: %w", err)
	}

	if opts.Hostname != "" {
		err := exec.Command("virt-customize", "-a", tempPath, "--hostname", opts.Hostname).Run()
		if err != nil {
			return fmt.Errorf("failed to set hostname via virt-customize: %w", err)
		}
		PrintVerbose(2, "Hostname set to: %s", opts.Hostname)
	}

	return nil
}
