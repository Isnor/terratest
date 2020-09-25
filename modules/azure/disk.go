package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// DiskExists indicates whether the specified Azure Managed Disk exists
// This function would fail the test if there is an error.
func DiskExists(t testing.TestingT, diskName string, resGroupName string, subscriptionID string) bool {
	exists, err := DiskExistsE(t, diskName, resGroupName, subscriptionID)
	require.NoError(t, err)
	return exists
}

// DiskExistsE indicates whether the specified Azure Managed Disk exists
func DiskExistsE(t testing.TestingT, diskName string, resGroupName string, subscriptionID string) (bool, error) {
	// Get the Disk object
	_, err := GetDiskE(t, diskName, resGroupName, subscriptionID)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetDiskType returns the Disk Storage Account Type of the Azure Managed Disk
// This property also accessible from the VM client disk storage object but only works
// when the VM is online, while this direct call to GetDiskType always works.
// This function would fail the test if there is an error.
func GetDiskType(t testing.TestingT, diskName string, resGroupName string, subscriptionID string) compute.DiskStorageAccountTypes {
	diskType, err := GetDiskTypeE(t, diskName, resGroupName, subscriptionID)
	require.NoError(t, err)
	return diskType
}

// GetDiskTypeE returns the Disk Storage Account Type of the Azure Managed Disk
func GetDiskTypeE(t testing.TestingT, diskName string, resGroupName string, subscriptionID string) (compute.DiskStorageAccountTypes, error) {
	// Get the Disk object
	disk, err := GetDiskE(t, diskName, resGroupName, subscriptionID)
	if err != nil {
		return "", err
	}

	return disk.Sku.Name, nil
}

// GetDisk returns a Disk in the specified Azure Resource Group
// This function would fail the test if there is an error.
func GetDisk(t testing.TestingT, diskName string, resGroupName string, subscriptionID string) *compute.Disk {
	disk, err := GetDiskE(t, diskName, resGroupName, subscriptionID)
	require.NoError(t, err)
	return disk
}

// GetDiskE returns a Disk in the specified Azure Resource Group
func GetDiskE(t testing.TestingT, diskName string, resGroupName string, subscriptionID string) (*compute.Disk, error) {
	// Validate resource group name and subscription ID
	resGroupName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	// Get the client reference
	client, err := GetDiskClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the Disk
	disk, err := client.Get(context.Background(), resGroupName, diskName)
	if err != nil {
		return nil, err
	}

	return &disk, nil
}

// GetDiskClientE returns a new Disk client in the specified Azure Subscription
func GetDiskClientE(subscriptionID string) (*compute.DisksClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the Disk client
	client := compute.NewDisksClient(subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer

	return &client, nil
}
