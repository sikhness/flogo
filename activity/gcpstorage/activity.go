package gcpstorage

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

const (
	operationWrite  = "WRITE"
	operationRead   = "READ"
	operationDelete = "DELETE"

	writeOptionNew       = "NEW"
	writeOptionOverwrite = "OVERWRITE"
	writeOptionAppend    = "APPEND"

	ivJSONCredentials = "jsonCredentials"
	ivBucketName      = "bucketName"
	ivOperation       = "operation"
	ivObjectName      = "objectName"
	ivObjectContent   = "objectContent"
	ivWriteOption     = "writeOption"

	ovOutput = "output"
)

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MyActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *MyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Function to login and create an authenticated client object to GCP Storage
func loginGCP(ctx context.Context, jsonCredentials string) (*storage.Client, error) {

	// Create a credentials object using provided GCP service account JSON Private Key
	creds, err := google.CredentialsFromJSON(ctx, []byte(jsonCredentials), storage.ScopeReadWrite)
	if err != nil {
		return nil, err
	}

	// Create an authenticated GCP Storage client to perform actions with
	client, err := storage.NewClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}
	_ = client // Use the client.

	return client, err
}

// Function to write text to a given object.  Information can be overwriten to an existing object and can
// also be appended to the given object
func writeObject(ctx context.Context, bkt *storage.BucketHandle, objectName string, objectContent string,
	writeOption string) (err error) {

	// Initialize the Object within the bucket. You can specify a folder structure as part of the
	// objectName as well
	obj := bkt.Object(objectName)

	// Initialize a new writer to the Object to prepare for writing
	w := obj.NewWriter(ctx)

	switch strings.ToUpper(writeOption) {
	case writeOptionNew:
		// Check to see if current Object exists by getting its contents
		currentObjectContent, err := readObject(ctx, bkt, objectName)

		// Return error if object already exists
		if currentObjectContent != "" && err == nil {
			return errors.New("Object already exists")
		}

		// Write text into the Object
		_, err = w.Write([]byte(objectContent))
		if err != nil {
			return err
		}
	case writeOptionAppend:
		// Read current Object object content (if exists)
		currentObjectContent, err := readObject(ctx, bkt, objectName)

		// Append text to current Object
		_, err = w.Write([]byte(currentObjectContent + objectContent))
		if err != nil {
			return err
		}
	case writeOptionOverwrite:
		// Overwrite text into the Object

		_, err = w.Write([]byte(objectContent))
		if err != nil {
			return err
		}
	default:
		err = errors.New("Unsupported write option")
		return err
	}

	// Close object after write is completed
	err = w.Close()
	if err != nil {
		return err
	}

	// Write succesful, no errors
	return nil
}

// Function to read text from a given object
func readObject(ctx context.Context, bkt *storage.BucketHandle, objectName string) (objectContent string, err error) {

	// Initialize the Object within the bucket. You can specify a folder structure as part of the
	// objectName as well
	obj := bkt.Object(objectName)

	// Initialize a new reader to the Object to prepare for reading
	r, err := obj.NewReader(ctx)
	if err != nil {
		return "", err
	}

	// Read the contents of the Object
	byteObjectContent, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}

	// Close the object object once done reading
	err = r.Close()
	if err != nil {
		return "", err
	}

	// Return object contents
	return string(byteObjectContent), nil
}

// Function to delete a given object in the bucket
func deleteObject(ctx context.Context, bkt *storage.BucketHandle, objectName string) (err error) {

	// Initialize the Object within the bucket. You can specify a folder structure as part of the
	// objectName as well
	obj := bkt.Object(objectName)

	// Delete the given object
	err = obj.Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Eval implements activity.Activity.Eval
func (a *MyActivity) Eval(ctx activity.Context) (done bool, err error) {

	// Gain activity inputs from context
	jsonCredentials, _ := ctx.GetInput(ivJSONCredentials).(string)
	bucketName, _ := ctx.GetInput(ivBucketName).(string)
	operation, _ := ctx.GetInput(ivOperation).(string)
	objectName, _ := ctx.GetInput(ivObjectName).(string)
	writeOption, _ := ctx.GetInput(ivWriteOption).(string)
	objectContent := fmt.Sprintf("%v", ctx.GetInput(ivObjectContent))

	gcpctx := context.Background()
	client, err := loginGCP(gcpctx, jsonCredentials)

	bkt := client.Bucket(bucketName)

	switch strings.ToUpper(operation) {
	case operationWrite:
		err = writeObject(gcpctx, bkt, objectName, objectContent, writeOption)
		if err != nil {
			return false, err
		}
	case operationRead:
		objectContent, err = readObject(gcpctx, bkt, objectName)
		if err != nil {
			return false, err
		}
		ctx.SetOutput(ovOutput, objectContent)
	case operationDelete:
		err = deleteObject(gcpctx, bkt, objectName)
		if err != nil {
			return false, err
		}
	default:
		err = errors.New("Unsupported operation")
		return false, err
	}

	// Evaluation complete, no errors
	return true, nil
}
