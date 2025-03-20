package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
)

// createBlobContainers creates a new Azure Blob Container within the specified storage account.
//
// Parameters:
//   - ctx: The context for the operation, used for cancellation and deadlines.
//   - resourceGroupName: The name of the resource group containing the storage account.
//   - storageAccountName: The name of the storage account where the blob container will be created.
//   - containerName: The name of the blob container to be created.
//
// Returns:
//   - A pointer to the created armstorage.BlobContainer object if successful.
//   - An error if the operation fails.
func createBlobContainers(ctx context.Context, resourceGroupName, storageAccountName, containerName string) (*armstorage.BlobContainer, error) {
	blobContainersClient, _ := BlobContainersClientClient()

	blobContainerResp, err := blobContainersClient.Create(
		ctx,
		resourceGroupName,
		storageAccountName,
		containerName,
		armstorage.BlobContainer{
			ContainerProperties: &armstorage.ContainerProperties{
				PublicAccess: to.Ptr(armstorage.PublicAccessNone),
			},
		},
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &blobContainerResp.BlobContainer, nil
}
