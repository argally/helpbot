package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
)

// checkNameAvailability checks the availability of a given storage account name in Azure.
// It uses the Storage Accounts Client to send a request to Azure and determine if the name
// is available for use.
//
// Parameters:
//   - ctx: The context for the request, which can include deadlines, timeouts, and other
//     cancellation signals.
//   - storageAccountName: The name of the storage account to check for availability.
//
// Returns:
//   - A pointer to an armstorage.CheckNameAvailabilityResult containing the result of the
//     name availability check.
//   - An error if the request to Azure fails or if there is an issue with the client.
func checkNameAvailability(ctx context.Context, storageAccountName string) (*armstorage.CheckNameAvailabilityResult, error) {
	saClient, _ := StorageAccountsClient()
	result, err := saClient.CheckNameAvailability(
		ctx,
		armstorage.AccountCheckNameAvailabilityParameters{
			Name: to.Ptr(storageAccountName),
			Type: to.Ptr("Microsoft.Storage/storageAccounts"),
		},
		nil)
	if err != nil {
		return nil, err
	}
	return &result.CheckNameAvailabilityResult, nil
}

// createStorageAccount creates an Azure Storage Account within the specified resource group.
// It allows specifying the storage account name, location, and uses default settings for encryption and access tier.
//
// Parameters:
//   - ctx: The context for the operation, used for cancellation and deadlines.
//   - resourceGroupName: The name of the resource group where the storage account will be created.
//   - storageAccountName: The name of the storage account to be created.
//   - location: The Azure region where the storage account will be created. If empty, a default location is used.
//
// Returns:
//   - *armstorage.Account: A pointer to the created storage account object.
//   - error: An error object if the operation fails, otherwise nil.
func createStorageAccount(ctx context.Context, resourceGroupName, storageAccountName, location string) (*armstorage.Account, error) {
	saClient, _ := StorageAccountsClient()
	// Set the location for the Resource Group
	if location == "" {
		location = DefaultLocation // Default location if none is provided
	}
	pollerResp, err := saClient.BeginCreate(
		ctx,
		resourceGroupName,
		storageAccountName,
		armstorage.AccountCreateParameters{
			Kind: to.Ptr(armstorage.KindStorageV2),
			SKU: &armstorage.SKU{
				Name: to.Ptr(armstorage.SKUNameStandardLRS),
			},
			Location: to.Ptr(location),
			Properties: &armstorage.AccountPropertiesCreateParameters{
				AccessTier: to.Ptr(armstorage.AccessTierCool),
				Encryption: &armstorage.Encryption{
					Services: &armstorage.EncryptionServices{
						File: &armstorage.EncryptionService{
							KeyType: to.Ptr(armstorage.KeyTypeAccount),
							Enabled: to.Ptr(true),
						},
						Blob: &armstorage.EncryptionService{
							KeyType: to.Ptr(armstorage.KeyTypeAccount),
							Enabled: to.Ptr(true),
						},
						Queue: &armstorage.EncryptionService{
							KeyType: to.Ptr(armstorage.KeyTypeAccount),
							Enabled: to.Ptr(true),
						},
						Table: &armstorage.EncryptionService{
							KeyType: to.Ptr(armstorage.KeyTypeAccount),
							Enabled: to.Ptr(true),
						},
					},
					KeySource: to.Ptr(armstorage.KeySourceMicrosoftStorage),
				},
			},
		}, nil)
	if err != nil {
		return nil, err
	}
	resp, err := pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &resp.Account, nil
}
