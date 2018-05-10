package gcpstorage

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
)

const (

	// Enter the json Private Key given when creating a service account in Google Cloud Platform
	jsonCredentials = `{
		<<PRIVATE KEY FROM GCP SERVICE ACCOUNT>>
	  }
	  `

	// Name of bucket pre-defined in GCP Storage
	bucketName = "<<BUCKET NAME>>"

	// Name of the object (file) that needs to be created in GCP Storage. Folder structures can be
	// defined here as well ex: "FOLDER/fileName.txt"
	objectName = "flogo.txt"

	/* Defines the ACL for the object being created.  Follow the JSON format defined below in the variable.
	More information about ACLs can be found at https://cloud.google.com/storage/docs/access-control/lists

	user: The user, group, project, or domain name to grant access. More information can be found in the link
	above. Example syntax for a user type would be "user-gcpuser@gmail.com"
	role: Can be either READER, WRITER, or OWNER
	*/
	objectACLList = `{
			"user-<<USER EMAIL>>": "OWNER",
			"user-<<USER EMAIL>>": "READER"
	  }`
)

var activityMetadata *activity.Metadata

func getActivityMetadata() *activity.Metadata {

	if activityMetadata == nil {
		jsonMetadataBytes, err := ioutil.ReadFile("activity.json")
		if err != nil {
			panic("No Json Metadata found for activity.json path")
		}

		activityMetadata = activity.NewMetadata(string(jsonMetadataBytes))
	}

	return activityMetadata
}

func TestCreate(t *testing.T) {

	act := NewActivity(getActivityMetadata())

	if act == nil {
		t.Error("Activity Not Created")
		t.Fail()
		return
	}
}

func TestCreateObject(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	tc.SetInput("jsonCredentials", jsonCredentials)
	tc.SetInput("bucketName", bucketName)
	tc.SetInput("operation", "WRITE")
	tc.SetInput("objectName", objectName)
	tc.SetInput("objectContent", "This text was input from the TestCreateObject test method\n")
	tc.SetInput("writeOption", "NEW")

	_, err := act.Eval(tc)
	if err != nil {
		panic(err)
	}
}

func TestCreateObjectOverwrite(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	tc.SetInput("jsonCredentials", jsonCredentials)
	tc.SetInput("bucketName", bucketName)
	tc.SetInput("operation", "WRITE")
	tc.SetInput("objectName", objectName)
	tc.SetInput("objectContent", "This text was input from the TestCreateObjectOverwrite test method\n")
	tc.SetInput("writeOption", "OVERWRITE")
	tc.SetInput("objectACLList", objectACLList)

	_, err := act.Eval(tc)
	if err != nil {
		panic(err)
	}
}

func TestCreateObjectAppend(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	tc.SetInput("jsonCredentials", jsonCredentials)
	tc.SetInput("bucketName", bucketName)
	tc.SetInput("operation", "WRITE")
	tc.SetInput("objectName", objectName)
	tc.SetInput("objectContent", "This text was input from the TestCreateObjectAppend test method\n")
	tc.SetInput("writeOption", "APPEND")
	tc.SetInput("objectACLList", objectACLList)

	_, err := act.Eval(tc)
	if err != nil {
		panic(err)
	}
}

func TestCreateObjectAppendNumber(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	tc.SetInput("jsonCredentials", jsonCredentials)
	tc.SetInput("bucketName", bucketName)
	tc.SetInput("operation", "WRITE")
	tc.SetInput("objectName", objectName)
	tc.SetInput("objectContent", 1234567890)
	tc.SetInput("writeOption", "APPEND")
	tc.SetInput("objectACLList", objectACLList)

	_, err := act.Eval(tc)
	if err != nil {
		panic(err)
	}
}

func TestReadObject(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	tc.SetInput("jsonCredentials", jsonCredentials)
	tc.SetInput("bucketName", bucketName)
	tc.SetInput("operation", "READ")
	tc.SetInput("objectName", objectName)

	_, err := act.Eval(tc)
	if err != nil {
		panic(err)
	}

	output := tc.GetOutput("output")
	fmt.Printf("Object Contents:\n%v\n\n", output)
}

func TestDeleteObject(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	tc.SetInput("jsonCredentials", jsonCredentials)
	tc.SetInput("bucketName", bucketName)
	tc.SetInput("operation", "DELETE")
	tc.SetInput("objectName", objectName)

	_, err := act.Eval(tc)
	if err != nil {
		panic(err)
	}

}
