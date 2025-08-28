package cmd

import (
	"os"
	"path/filepath"
	"reflect"
)

type Options struct {
	Profile             string `yaml:"-"`                           // Internal use only, not from YAML
	VirtMode            string `yaml:"virt_mode"`                   // "install" or "customize"
	ImageName           string `yaml:"image_name"`                  // Name of the image
	OSVariant           string `yaml:"os_var"`                      // OS variant (from osinfo-query)
	ImageSize           string `yaml:"image_size"`                  // e.g., "100G"
	ImagePartition 		string `yaml:"image_part"`					// e.g. "/dev/sda1"
	KickstartPath       string `yaml:"ks_file" file:"true"`         // Path to Kickstart file
	ISOFile             string `yaml:"iso_file" file:"true"`        // Local ISO file path
	RepoURL             string `yaml:"repo_url"`                    // Remote repo URL
	ImageSource         string `yaml:"image_src" file:"true"`       // Source QCOW2 image (customize mode)
	ImageDestination    string `yaml:"image_dest" file:"true"`      // Output QCOW2 image
	Hostname            string `yaml:"hostname"`                    // Optional
	CustomScript        string `yaml:"custom_script" file:"true"`   // Optional bash script
	Network             string `yaml:"network"`                     // Optional virtual network
	Console             string `yaml:"console"`                     // "serial", or "graphical"
	Firmware            string `yaml:"firmware"`                    // "bios" or "efi" (default: bios)
}

func resolvePath(path string) string {
	if path == "" || filepath.IsAbs(path) {
		return path
	}
	cwd, err := os.Getwd()
	if err != nil {
		return path
	}
	return filepath.Join(cwd, path)
}

func (o *Options) ResolvePaths() {
	v := reflect.ValueOf(o).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("file")
		if tag == "true" && field.Type.Kind() == reflect.String {
			val := v.Field(i).String()
			v.Field(i).SetString(resolvePath(val))
		}
	}
}
