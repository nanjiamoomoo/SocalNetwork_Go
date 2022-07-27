package backend

// [START storage_upload_file]
import (
	"context"
	"fmt"
	"io"
	"os"

	"socialnetwork_go/constants"

	"cloud.google.com/go/storage"
)

var (
	GCSBackend *GoogleCloudStorageBackend
)

type GoogleCloudStorageBackend struct {
	client *storage.Client
	bucket string
}

//obtain a new client for GCS
func InitGCSBackend() {
	//use cloud service account key to access to CGS 
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "C:/Users/ZW/Downloads/socialnetwork-357505-594450545c97.json")
    client, err := storage.NewClient(context.Background())
    if err != nil {
        panic(err)
    }

    GCSBackend = &GoogleCloudStorageBackend{
        client: client,
        bucket: constants.GCS_BUCKET,
    }
}

//save file to GCS and return the MediaLink of the file
func (backend *GoogleCloudStorageBackend) SaveToGCS(r io.Reader, objectName string) (string, error) {
	ctx := context.Background()

	//Bucket returns a BucketHandle, which provides operations on the named bucket.
	//Object returns an ObjectHandle, which provides operations on the named object.
	object := backend.client.Bucket(backend.bucket).Object(objectName)

	//NewWriter returns a storage Writer that writes to the GCS object
	wc := object.NewWriter(ctx)

	// Copy copies from reader(src) to writer(dst) until either EOF is reached
	// on src or an error occurs. It returns the number of bytes
	// copied and the first error encountered while copying, if any.
	if _, err := io.Copy(wc, r); err != nil {
		return "", err
	}

	if err := wc.Close(); err != nil {
		return "", err
	}

	//ACL provides access to the object's access control list.
	//set all users can read the file
	if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", err
	}

	//Attrs returns meta information about the object, such us URL
	attrs, err := object.Attrs(ctx)
	if err != nil {
		return "", err
	}

	fmt.Printf("File is saved to GCS: %s\n", attrs.MediaLink)
	return attrs.MediaLink, nil
}
