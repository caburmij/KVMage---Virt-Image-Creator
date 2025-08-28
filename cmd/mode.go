package cmd

import "errors"

func ValidateModeFlags(runMode bool, configPath string) error {
	switch {
	case runMode && configPath != "":
		return errors.New("cannot specify both --run and --config")
	case !runMode && configPath == "":
		return errors.New("must specify either --run or --config")
	default:
		return nil
	}
}
