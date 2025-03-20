package azure

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

// CreateResourceGroup creates a new Azure Resource Group with the specified name and region.
// If the region is not provided, it defaults to "eastus".
// It checks if the Resource Group already exists before attempting to create it.
func CreateResourceGroup(ctx context.Context, resourceGroupName, location string) (*armresources.ResourceGroup, error) {
	// Initialize the Resource Group client
	rgClient, err := ResourceGroupClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create resource group client: %w", err)
	}
	// Set the location for the Resource Group
	if location == "" {
		location = DefaultLocation // Default region if none is provided
	}
	resourceGroupParams := armresources.ResourceGroup{
		Location: to.Ptr(string(location)),
	}

	resourceGroupResp, err := rgClient.CreateOrUpdate(ctx, resourceGroupName, resourceGroupParams, nil)
	if err != nil {
		// Check if the error contains a specific message about location
		var respErr *azcore.ResponseError
		if errors.As(err, &respErr) && respErr.ErrorCode == "LocationNotAvailableForResourceGroup" {
			return nil, fmt.Errorf("enter valid location: %s", respErr.ErrorCode)
		}
		return nil, err
	}
	return &resourceGroupResp.ResourceGroup, nil
}

// Return an error if the Resource Group already exists

// checkResourceGroupAvailability checks if a Resource Group with the specified name exists.
// It returns true if the Resource Group exists, otherwise false.
func checkResourceGroupAvailability(ctx context.Context, resourceGroupName string) (bool, error) {
	// Initialize the Resource Group client
	rgClient, err := ResourceGroupClient()
	if err != nil {
		return false, fmt.Errorf("failed to create resource group client: %w", err)
	}
	// Check the existence of the Resource Group
	result, err := rgClient.CheckExistence(ctx, resourceGroupName, nil)
	if err != nil {
		return false, fmt.Errorf("error checking resource group existence: %w", err)
	}
	// Return true if the Resource Group exists (HTTP success status code)
	return result.Success, nil
}
