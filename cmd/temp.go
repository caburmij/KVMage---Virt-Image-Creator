package cmd

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	TempImageName     string
	TempImagePath     string
	TempKickstartPath string
	TempInstallMedia  string
	TempImageSource   string
	TempCustomScript  string
)

const tempDir = "/var/lib/libvirt/images"

func CreateTempImage(opts *Options) (string, string, error) {
	RequireRoot()

	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", "", fmt.Errorf("generate random name: %w", err)
	}
	TempImageName = fmt.Sprintf("kvmage-%s", hex.EncodeToString(b))
	TempImagePath = filepath.Join(tempDir, TempImageName+".qcow2")

	if opts.VirtMode == "install" {
		PrintVerbose(2, "Temporary image path: %s", TempImagePath)
		PrintVerbose(2, "Requested image size: %s", opts.ImageSize)

		args := []string{"create", "-f", "qcow2", "-o", "compat=0.10", TempImagePath, opts.ImageSize}
		PrintVerbose(3, "Running command: qemu-img %s", joinArgs(args))

		err := exec.Command("qemu-img", args...).Run()
		if err != nil {
			return "", "", fmt.Errorf("qemu-img create failed: %w", err)
		}
	} else if opts.VirtMode == "customize" {
		PrintVerbose(2, "Copying source image %s to %s", opts.ImageSource, TempImagePath)

		args := []string{"convert", "-O", "qcow2", opts.ImageSource, TempImagePath}
		PrintVerbose(3, "Running command: qemu-img %s", joinArgs(args))

		err := exec.Command("qemu-img", args...).Run()
		if err != nil {
			return "", "", fmt.Errorf("qemu-img convert failed: %w", err)
		}
	} else {
		return "", "", fmt.Errorf("invalid VirtMode: must be 'install' or 'customize'")
	}

	return TempImageName, TempImagePath, nil
}

func resolveRemoteISO(src string) (string, error) {
	if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
		filename := filepath.Base(src)
		dest := filepath.Join(os.TempDir(), filename)

		if _, err := os.Stat(dest); err == nil {
			PrintVerbose(2, "Using cached remote ISO: %s", dest)
			return dest, nil
		}

		var downloader string
		if _, err := exec.LookPath("curl"); err == nil {
			downloader = "curl"
		} else if _, err := exec.LookPath("wget"); err == nil {
			downloader = "wget"
		} else {
			return "", fmt.Errorf("neither curl nor wget is installed")
		}

		PrintVerbose(2, "Downloading ISO from %s using %s", src, downloader)

		var cmd *exec.Cmd
		if downloader == "curl" {
			cmd = exec.Command("curl", "-L", "-v", "-o", dest, src)
		} else {
			cmd = exec.Command("wget", "-v", "-O", dest, src)
		}

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("failed to download remote ISO: %w", err)
		}

		return dest, nil
	}

	return src, nil
}

func copyToTemp(src string, label string) (string, error) {
	if src == "" {
		return "", nil
	}

	ext := filepath.Ext(src)
	destName := fmt.Sprintf("%s-%s.temp%s", TempImageName, label, ext)
	dest := filepath.Join(tempDir, destName)

	PrintVerbose(2, "Copying %s to temp file: %s", label, dest)

	in, err := os.Open(src)
	if err != nil {
		return "", fmt.Errorf("failed to open %s: %w", src, err)
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file %s: %w", dest, err)
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return "", fmt.Errorf("failed to copy to temp file %s: %w", dest, err)
	}

	return dest, nil
}

func CopyInputFilesToTempDir(opts *Options) error {
	RequireRoot()

	var err error

	if TempKickstartPath, err = copyToTemp(opts.KickstartPath, "ks"); err != nil {
		return fmt.Errorf("kickstart file copy failed: %w", err)
	}

	if opts.ISOFile != "" {
		var resolved string
		if resolved, err = resolveRemoteISO(opts.ISOFile); err != nil {
			return fmt.Errorf("ISO resolution failed: %w", err)
		}
		if TempInstallMedia, err = copyToTemp(resolved, "iso"); err != nil {
			return fmt.Errorf("iso file copy failed: %w", err)
		}
	} else if opts.RepoURL != "" {
		TempInstallMedia = opts.RepoURL
		PrintVerbose(2, "Using repo URL as install media: %s", TempInstallMedia)
	}

	if TempImageSource, err = copyToTemp(opts.ImageSource, "src"); err != nil {
		return fmt.Errorf("source image copy failed: %w", err)
	}
	if TempCustomScript, err = copyToTemp(opts.CustomScript, "custom"); err != nil {
		return fmt.Errorf("custom script copy failed: %w", err)
	}

	return nil
}

func joinArgs(args []string) string {
	return fmt.Sprintf("%q", args)
}
