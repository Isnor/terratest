package terraform_test

import (
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalTfArgs(t *testing.T) {

	type MarshalTfTest struct {
		name            string
		args            *ExampleCommandArgs
		expectedCommand string
	}

	for _, test := range []MarshalTfTest{
		{
			name:            "migrate set",
			expectedCommand: "-migrate-state -force-copy",
			args: &ExampleCommandArgs{
				MigrateState: true,
				Refresh:      true,
			},
		},
		{
			name:            "plugin dir set",
			expectedCommand: "-plugin-dir=/home/foobar/plugins",
			args: &ExampleCommandArgs{
				PluginDir: "/home/foobar/plugins",
				Refresh:   true,
			},
		},
		{
			name:            "backend config set",
			expectedCommand: "-backend-config=foo=bar -backend-config=prefix.foo=prefix.bar -backend-config=zoobar={\"mymaparg\" = null}",
			args: &ExampleCommandArgs{
				BackendConfig: map[string]any{
					"foo":        "bar",
					"prefix.foo": "prefix.bar",
					"zoobar":     map[string]any{"mymaparg": nil},
				},
				Refresh: true,
			},
		},
		{
			name:            "var file set",
			expectedCommand: "-var-file foo -var-file bar",
			args: &ExampleCommandArgs{
				VarFiles: []string{"foo", "bar"},
				Refresh:  true,
			},
		},
	} {
		t.Run(test.name, func(subtest *testing.T) {
			cmdString, err := terraform.MarshalTfArgs(test.args)
			assert.NoError(subtest, err)
			assert.Equal(subtest, test.expectedCommand, cmdString)
		})
	}
}

func TestMarshalOptionsSuccess(t *testing.T) {

	type MarshalOptionsTest struct {
		name         string
		options      *terraform.Options
		expectations func(parent *testing.T, expectedCommand string, err error)
	}

	for _, test := range []MarshalOptionsTest{
		{
			name: "init success",
			options: &terraform.Options{
				TerraformBinary: "terraform",
				NoColor:         true,
				AutoApprove:     true,
				MigrateState:    true,
				Upgrade:         true,
			},
			expectations: func(parent *testing.T, expectedCommand string, err error) {
				require.NoError(parent, err)
				// TODO: right now we're ignoring the actual command, this is just for the arguments
				assert.ElementsMatch(parent, strings.Split(expectedCommand, " "), []string{"-upgrade", "-migrate-state", "-force-copy", "-auto-approve", "-no-color", "-refresh=false", "-input=false"})
			},
		},
	} {
		t.Run(test.name, func(subtest *testing.T) {
			cmdString, err := terraform.MarshalTfArgs(test.options)
			test.expectations(subtest, cmdString, err)
		})
	}
}

// ExampleCommandArgs are command argument examples of different types of arguments from the existing commands
type ExampleCommandArgs struct {
	BackendConfig map[string]any `tfarg:"-backend-config"`                      // The vars to pass to the terraform init command for extra configuration for the backend
	MigrateState  bool           `tfarg:"-migrate-state -force-copy,omitempty"` // Set the -migrate-state and -force-copy (suppress 'yes' answer prompt) flag to the terraform init command
	PluginDir     string         `tfarg:"-plugin-dir"`                          // The path of downloaded plugins to pass to the terraform init command (-plugin-dir)
	VarFiles      []string       `tfarg:"-var-file"`                            // The var file paths to pass to Terraform commands using -var-file option.
	// reconfigure is also a bool type, but it's encoded slightly differently than MigrateState
	Reconfigure bool `tfarg:"-reconfigure,omitempty"` // Set the -reconfigure flag to the terraform init command
	// refresh is another odd bool case where the default is actually true, so when this is supplied as false we need to write `-refresh=false`
	Refresh bool `tfarg:"-refresh,omittrue"`
}
