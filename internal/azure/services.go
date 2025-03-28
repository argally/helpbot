package azure

import (
	"context"
	"fmt"
	"log"

	utils "github.com/argally/helpbot/internal/util"
)

// AzureServicesCreate creates an Azure resource based on the specified parameters.
// It supports creating resource groups, storage accounts, and blob containers.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellations.
//   - resourceName: The base name for the resource to be created.
//   - resourceType: The type of resource to create. Valid values are "storage-account" and "blob-container".
//   - location: The Azure region where the resource will be created.
//
// Returns:
//   - A string message indicating the result of the operation.
//   - An error if the operation fails.
//
// Behavior:
//   - If the resource group does not exist, it will be created.
//   - If the resourceType is "storage-account", a storage account will be created within the resource group.
//   - If the resourceType is "blob-container", both a storage account and a blob container will be created.
//   - If the resourceType is invalid, an error will be returned.
//
// Notes:
//   - The storage account and blob container names are generated by appending the current date to the resourceName.
//   - If a storage account name is unavailable, a warning message will be returned instead of creating the resource.
func AzureServicesCreate(ctx context.Context, resourceName, resourceType, location string) (string, error) {
	resourceGroupName := utils.ConcatStrings(resourceName, "-rg")
	resourceGroupExists, err := checkResourceGroupAvailability(ctx, resourceGroupName)
	if err != nil {
		return "", err
	}
	if resourceGroupExists {
		log.Printf("Resource Group already exists will re-use: %s", resourceGroupName)
	}
	resourceGroupResp, err := CreateResourceGroup(ctx, resourceGroupName, location)
	if err != nil {
		return "", err
	}
	log.Printf("Creating Resource Group: %v", *resourceGroupResp.Name)
	// If the resource type is "storageAccount", ensure the storage account is created
	if resourceType == "storage-account" || resourceType == "blob-container" {
		storageAccountName := utils.ConcatStrings(resourceName, utils.GetCurrentDateFormatted())

		// Check if the storage account name is available
		storageAccountAvailable, err := checkNameAvailability(ctx, storageAccountName)
		if err != nil {
			return "", err
		}
		if !*storageAccountAvailable.NameAvailable {
			log.Printf("Storage account already exists duplicate: %v", *storageAccountAvailable.Message)
			return fmt.Sprintf(":warning: %s", *storageAccountAvailable.Message), nil
		}

		// Create the storage account
		storageAccountResp, err := createStorageAccount(ctx, resourceGroupName, storageAccountName, location)
		if err != nil {
			return "", err
		}
		log.Printf("Creating Storage account: %v", *storageAccountResp.Name)
		if resourceType == "storage-account" {
			return fmt.Sprintf(":white_check_mark: Storage account '%s' created successfully in resource group '%s'", *storageAccountResp.Name, *resourceGroupResp.Name), nil
		}
		blobContainerName := utils.ConcatStrings(resourceName, utils.GetCurrentDateFormatted())
		// Create the storage account
		blobContainreResp, err := createBlobContainers(ctx, resourceGroupName, storageAccountName, blobContainerName)
		if err != nil {
			return "", err
		}
		log.Printf("Creating Blob Container: %v", *blobContainreResp.Name)
		return fmt.Sprintf(":white_check_mark: Blob Container '%s' and Storage account '%s' created successfully in resource group '%s'", *blobContainreResp.Name, *storageAccountResp.Name, *resourceGroupResp.Name), nil
	}

	// If an invalid resource type is provided, return an error
	return "", fmt.Errorf("invalid resource type: %s. Valid types are 'blob-container' or 'storage-account'", resourceType)
}
