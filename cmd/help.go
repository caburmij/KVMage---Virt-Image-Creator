package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func CustomHelp(cmd *cobra.Command, args []string) {
	Print("KVMage - Build and customize KVM images using virt-install and virt-customize.\n")
	Print("")
	Print("Usage:")
	Print("  kvmage [--run | --config] [flags]")
	Print("")
	Print("Execution Modes (required):")
	Print("  -r, --run                     Use CLI arguments directly")
	Print("  -f, --config <file>           Use a YAML config file")
	Print("")
	Print("Installation Methods (required):")
	Print("  -i, --install                 Install mode (create image from ISO)")
	Print("  -c, --customize               Customize mode (modify existing image)")
	Print("")
	Print("Image Options:")
	Print("  -n, --image-name <name>       Name of the image")
	Print("  -o, --os-var <os>             OS variant (use `osinfo-query os`)")
	Print("  -s, --image-size <size>       Image size (e.g., 100G), expands image in customize mode")
	Print("  -P, --image-part <device>     Partition to expand (e.g., /dev/sda1)")
	Print("  -k, --ks-file <file>          Path to Kickstart file")
	Print("  -l, --install-media <path>    Install media path or URL")
	Print("  -S, --image-src <file>        Source QCOW2 image (customize mode)")
	Print("  -D, --image-dest <file>       Destination QCOW2 image")
	Print("  -H, --hostname <name>         Hostname to set inside the image (optional)")
	Print("  -C, --custom-script <file>    Bash script to run inside image (optional)")
	Print("  -W, --network <iface>         Virtual network name (optional)")
	Print("  -m, --firmware <type>         Firmware type: bios (default) or efi")
	Print("")
	Print("Global Options:")
	Print("  -h, --help                    Show help and exit")
	Print("  -v, --verbose                 Enable verbose output (-v/-vv/-vvv)")
	Print("      --verbose-level <n>       Set verbosity level explicitly (0-3)")
	Print("  -q, --quiet                   Suppress all output")
	Print("  -V, --version                 Show version info for KVMage and tools")
	Print("  -u, --uninstall               Uninstall KVMage from /usr/local/bin")
	Print("  -X, --cleanup                 Run cleanup mode to remove orphaned kvmage temp files")
	Print("")
	os.Exit(0)
}
