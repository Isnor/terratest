// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package azure

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/require"
)

func TestGetTargetAzureSubscription(t *testing.T) {
	t.Parallel()

	var id string
	var exists bool

	//Lookup ARM_SUBSCRIPTION_ID env variable, CI requires this value to run all test.
	if id, exists = os.LookupEnv(azure.AzureSubscriptionID); !exists {
		fmt.Printf("ARM_SUBSCRIPTION_ID environment variable not set")
		id = ""
	}

	type args struct {
		subID string
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "subIDProvidedAsArg", args: args{subID: "test"}, want: "test", wantErr: false},
		{name: "subIDNotProvidedFallbackToEnv", args: args{subID: ""}, want: id, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := azure.GetTargetAzureSubscription(tt.args.subID)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestGetTargetAzureResourceGroupName(t *testing.T) {
	t.Parallel()

	type args struct {
		rgName string
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "rgNameProvidedAsArg", args: args{rgName: "test"}, want: "test", wantErr: false},
		{name: "rgNameNotProvided", args: args{rgName: ""}, want: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := azure.GetTargetAzureResourceGroupName(tt.args.rgName)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.Equal(t, tt.want, got)
			}
		})
	}
}
