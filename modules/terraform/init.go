package terraform

import (
	"github.com/gruntwork-io/terratest/modules/testing"
)

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
	// args := []string{"init", fmt.Sprintf("-upgrade=%t", options.Upgrade || options.InitArgs.Upgrade)}

	// // all of the deprecated arguments from Options that were passed to Init need to be declared here
	// // and retrieved from either Options or Options.InitArgs below
	// var backendConfig map[string]any
	// var pluginDir string

	// if len(options.BackendConfig) != 0 {
	// 	backendConfig = options.BackendConfig
	// } else {
	// 	backendConfig = options.InitArgs.BackendConfig
	// }
	// if len(options.PluginDir) > 0 {
	// 	pluginDir = options.PluginDir
	// } else {
	// 	pluginDir = options.InitArgs.PluginDir
	// }

	// // append the args, from wherever they came from (Options or Options.InitArgs)
	// if options.Reconfigure || options.InitArgs.Reconfigure {
	// 	args = append(args, "-reconfigure")
	// }
	// if options.MigrateState || options.InitArgs.MigrateState {
	// 	args = append(args, "-migrate-state", "-force-copy")

	// }
	// // Append no-color option if needed
	// if options.NoColor {
	// 	args = append(args, "-no-color")
	// }

	// args = append(args, FormatTerraformBackendConfigAsArgs(backendConfig)...)
	// args = append(args, FormatTerraformPluginDirAsArgs(pluginDir)...)
	args, err := MarshalTfArgs(options)
	if err != nil {
		return "", err
	}
	return RunTerraformCommandE(t, options, "init", args)
}
