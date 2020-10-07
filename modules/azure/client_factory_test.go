/*

This file contains unit tests for the client factory implementation(s).

*/

package azure

import (
	"os"
	"testing"

	autorest "github.com/Azure/go-autorest/autorest/azure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Local consts for this file only
const govCloudEnvName = "AzureUSGovernmentCloud"
const publicCloudEnvName = "AzurePublicCloud"
const chinaCloudEnvName = "AzureChinaCloud"
const germanyCloudEnvName = "AzureGermanCloud"

func TestDefaultEnvIsPublicWhenNotSet(t *testing.T) {
	// save any current env value and restore on exit
	originalEnv := os.Getenv(AzureEnvironmentEnvName)
	defer os.Setenv(AzureEnvironmentEnvName, originalEnv)

	// Set env var to missing value
	os.Setenv(AzureEnvironmentEnvName, "")

	// get the default
	env := getDefaultEnvironmentName()

	// Make sure it's public cloud
	assert.Equal(t, autorest.PublicCloud.Name, env)
}

func TestDefaultEnvSetToGov(t *testing.T) {
	// save any current env value and restore on exit
	originalEnv := os.Getenv(AzureEnvironmentEnvName)
	defer os.Setenv(AzureEnvironmentEnvName, originalEnv)

	// Set env var to gov
	os.Setenv(AzureEnvironmentEnvName, govCloudEnvName)

	// get the default
	env := getDefaultEnvironmentName()

	// Make sure it's public cloud
	assert.Equal(t, autorest.USGovernmentCloud.Name, env)
}

func TestSubscriptionClientBaseURISetCorrectly(t *testing.T) {
	var cases = []struct {
		CaseName        string
		EnvironmentName string
		ExpectedBaseURI string
	}{
		{"GovCloud/SubscriptionClient", govCloudEnvName, autorest.USGovernmentCloud.ResourceManagerEndpoint},
		{"PublicCloud/SubscriptionClient", publicCloudEnvName, autorest.PublicCloud.ResourceManagerEndpoint},
		{"ChinaCloud/SubscriptionClient", chinaCloudEnvName, autorest.ChinaCloud.ResourceManagerEndpoint},
		{"GermanCloud/SubscriptionClient", germanyCloudEnvName, autorest.GermanCloud.ResourceManagerEndpoint},
	}

	// save any current env value and restore on exit
	currentEnv := os.Getenv(AzureEnvironmentEnvName)
	defer os.Setenv(AzureEnvironmentEnvName, currentEnv)

	for _, tt := range cases {
		t.Run(tt.CaseName, func(t *testing.T) {
			// Override env setting
			os.Setenv(AzureEnvironmentEnvName, tt.EnvironmentName)

			// Get a VM client
			client, err := GetClientForSubscriptionsE()
			require.NoError(t, err)

			// Check for correct ARM URI
			assert.Equal(t, tt.ExpectedBaseURI, client.BaseURI)
		})
	}
}

func TestVMClientBaseURISetCorrectly(t *testing.T) {
	var cases = []struct {
		CaseName        string
		EnvironmentName string
		ExpectedBaseURI string
	}{
		{"GovCloud/VMClient", govCloudEnvName, autorest.USGovernmentCloud.ResourceManagerEndpoint},
		{"PublicCloud/VMClient", publicCloudEnvName, autorest.PublicCloud.ResourceManagerEndpoint},
		{"ChinaCloud/VMClient", chinaCloudEnvName, autorest.ChinaCloud.ResourceManagerEndpoint},
		{"GermanCloud/VMClient", germanyCloudEnvName, autorest.GermanCloud.ResourceManagerEndpoint},
	}

	// save any current env value and restore on exit
	currentEnv := os.Getenv(AzureEnvironmentEnvName)
	defer os.Setenv(AzureEnvironmentEnvName, currentEnv)

	for _, tt := range cases {
		t.Run(tt.CaseName, func(t *testing.T) {
			// Override env setting
			os.Setenv(AzureEnvironmentEnvName, tt.EnvironmentName)

			// Get a VM client
			client, err := GetClientForVirtualMachinesE("")
			require.NoError(t, err)

			// Check for correct ARM URI
			assert.Equal(t, tt.ExpectedBaseURI, client.BaseURI)
		})
	}
}

func TestManagedClustersClientBaseURISetCorrectly(t *testing.T) {
	var cases = []struct {
		CaseName        string
		EnvironmentName string
		ExpectedBaseURI string
	}{
		{"GovCloud/ManagedClustersClient", govCloudEnvName, autorest.USGovernmentCloud.ResourceManagerEndpoint},
		{"PublicCloud/ManagedClustersClient", publicCloudEnvName, autorest.PublicCloud.ResourceManagerEndpoint},
		{"ChinaCloud/ManagedClustersClient", chinaCloudEnvName, autorest.ChinaCloud.ResourceManagerEndpoint},
		{"GermanCloud/ManagedClustersClient", germanyCloudEnvName, autorest.GermanCloud.ResourceManagerEndpoint},
	}

	// save any current env value and restore on exit
	currentEnv := os.Getenv(AzureEnvironmentEnvName)
	defer os.Setenv(AzureEnvironmentEnvName, currentEnv)

	for _, tt := range cases {
		t.Run(tt.CaseName, func(t *testing.T) {
			// Override env setting
			os.Setenv(AzureEnvironmentEnvName, tt.EnvironmentName)

			// Get a VM client
			client, err := GetClientForManagedClustersE("")
			require.NoError(t, err)

			// Check for correct ARM URI
			assert.Equal(t, tt.ExpectedBaseURI, client.BaseURI)
		})
	}
}
