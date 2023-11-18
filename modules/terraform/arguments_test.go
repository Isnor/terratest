package terraform_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestMarshalTfArgs(t *testing.T) {

	type MarshalTfTest struct {
		name            string
		args            *terraform.ExampleCommandArgs
		expectedCommand string
	}

	for _, test := range []MarshalTfTest{
		{
			name:            "migrate set",
			expectedCommand: "-migrate-state -force-copy",
			args: &terraform.ExampleCommandArgs{
				MigrateState: true,
				Refresh:      true,
			},
		},
		{
			name:            "plugin dir set",
			expectedCommand: "-plugin-dir=/home/foobar/plugins",
			args: &terraform.ExampleCommandArgs{
				PluginDir: "/home/foobar/plugins",
				Refresh:   true,
			},
		},
		{
			name:            "backend config set",
			expectedCommand: "-backend-config=foo=bar -backend-config=prefix.foo=prefix.bar",
			args: &terraform.ExampleCommandArgs{
				BackendConfig: map[string]any{
					"foo": "bar",
					"prefix.foo": "prefix.bar",
				},
				Refresh: true,
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			cmdString, err := terraform.MarshalTfArgs(test.args)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedCommand, cmdString)
		})
	}
}
