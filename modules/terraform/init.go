package terraform

import (
	"fmt"

	"github.com/gruntwork-io/terratest/modules/testing"
)

type InitArgs struct {
	BackendConfig map[string]any // The vars to pass to the terraform init command for extra configuration for the backend
	Upgrade       bool           // Whether the -upgrade flag of the terraform init command should be set to true or not
	Reconfigure   bool           // Set the -reconfigure flag to the terraform init command
	MigrateState  bool           // Set the -migrate-state and -force-copy (suppress 'yes' answer prompt) flag to the terraform init command
	PluginDir     string         // The path of downloaded plugins to pass to the terraform init command (-plugin-dir)
}

// Init calls terraform init and return stdout/stderr.
func Init(t testing.TestingT, options *Options) string {
	out, err := InitE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// InitE calls terraform init and return stdout/stderr.
func InitE(t testing.TestingT, options *Options) (string, error) {
	args := []string{"init", fmt.Sprintf("-upgrade=%t", options.Upgrade || options.InitArgs.Upgrade)}

	// all of the deprecated arguments from Options that were passed to Init need to be declared here
	// and retrieved from either Options or Options.InitArgs below
	var backendConfig map[string]any
	var pluginDir string

	if len(options.BackendConfig) != 0 {
		backendConfig = options.BackendConfig
	} else {
		backendConfig = options.InitArgs.BackendConfig
	}
	if len(options.PluginDir) > 0 {
		pluginDir = options.PluginDir
	} else {
		pluginDir = options.InitArgs.PluginDir
	}

	// append the args, from wherever they came from (Options or Options.InitArgs)
	if options.Reconfigure || options.InitArgs.Reconfigure {
		args = append(args, "-reconfigure")
	}
	if options.MigrateState || options.InitArgs.MigrateState {
		args = append(args, "-migrate-state", "-force-copy")

	}
	// Append no-color option if needed
	if options.NoColor {
		args = append(args, "-no-color")
	}

	args = append(args, FormatTerraformBackendConfigAsArgs(backendConfig)...)
	args = append(args, FormatTerraformPluginDirAsArgs(pluginDir)...)
	return RunTerraformCommandE(t, options, args...)
}

// each command has a struct with fields for each argument it supports, which we get from terraform <cmd> -help
// in the "CmdE" function for each command, we look at the fields that are set and use the tag of the field to
// build the arg list string
