package services

import (
	"context"
	"net/url"
	"os"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
	uuid "github.com/satori/go.uuid"
)

func UploadImageToCandidateAzureStorage(img []byte) (string, error) {
	key, account, container, endpoint := accountInfo()
	u, _ := url.Parse(endpoint + container + "/" + blobName())
	credential, err := azblob.NewSharedKeyCredential(account, key)
	if err != nil {
		return "", err
	}

	blockBlobUrl := azblob.NewBlockBlobURL(*u, azblob.NewPipeline(credential, azblob.PipelineOptions{}))

	cxt := context.Background()

	options := azblob.UploadToBlockBlobOptions{
		BlobHTTPHeaders: azblob.BlobHTTPHeaders{
			ContentType: "image/jpg",
		},
	}

	_, err = azblob.UploadBufferToBlockBlob(cxt, img, blockBlobUrl, options)
	if err != nil {
		return "", err
	}

	return blockBlobUrl.String(), nil
}

func accountInfo() (string, string, string, string) {
	return os.Getenv("AZURE_ACCESS_KEY"), os.Getenv("AZURE_BLOB_ACCOUNT_NAME"), os.Getenv("AZURE_BLOB_CONTAINER_NAME"), os.Getenv("AZURE_BLOB_SERVICE_ENDPOINT")
}

func blobName() string {
	t := time.Now()
	uuid := uuid.NewV4()

	return t.Format("20060102") + "-" + uuid.String() + ".jpg"
}
