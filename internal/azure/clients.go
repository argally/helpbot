package azure

import (
	"log"
	"os"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/joho/godotenv"
)

var (
	_                         = godotenv.Load(".env")
	onceAzureCred             sync.Once
	onceResourceGroupClient   sync.Once
	onceStorageAccountsClient sync.Once
	onceBlobContainersClient  sync.Once
	azureCred                 *azidentity.DefaultAzureCredential
	resourceGroupClient       *armresources.ResourceGroupsClient
	storageAccountsClient     *armstorage.AccountsClient
	blobContainersClient      *armstorage.BlobContainersClient
)

func AzureCredential() (*azidentity.DefaultAzureCredential, error) {
	// Create a new service client with token credential
	onceAzureCred.Do(func() {
		var err error

		azureCred, err = azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			log.Printf("unable to load SDK config, %v", err)
		}
	})

	return azureCred, nil
}

func ResourceGroupClient() (*armresources.ResourceGroupsClient, error) {
	onceResourceGroupClient.Do(func() {
		azureCred, _ = AzureCredential()
		var err error

		resourceGroupClient, err = armresources.NewResourceGroupsClient(os.Getenv("AZURE_SUBSCRIPTION_ID"), azureCred, nil)
		if err != nil {
			log.Printf("Issue connecting to Resource group client API, %v", err)
		}
	})

	return resourceGroupClient, nil
}

func StorageAccountsClient() (*armstorage.AccountsClient, error) {
	onceStorageAccountsClient.Do(func() {
		azureCred, _ = AzureCredential()
		var err error

		storageAccountsClient, err = armstorage.NewAccountsClient(os.Getenv("AZURE_SUBSCRIPTION_ID"), azureCred, nil)
		if err != nil {
			log.Printf("Issue connecting to Storage Account client API, %v", err)
		}
	})

	return storageAccountsClient, nil
}

func BlobContainersClientClient() (*armstorage.BlobContainersClient, error) {
	onceBlobContainersClient.Do(func() {
		azureCred, _ = AzureCredential()
		var err error

		blobContainersClient, err = armstorage.NewBlobContainersClient(os.Getenv("AZURE_SUBSCRIPTION_ID"), azureCred, nil)
		if err != nil {
			log.Printf("Issue connecting to Blob Container client API, %v", err)
		}
	})

	return blobContainersClient, nil
}
