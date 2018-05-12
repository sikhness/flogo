# GCPStorage
This activity allows you to interact with objects within a bucket in Google Cloud Platform - Storage

## Installation
### Flogo CLI  
```bash
flogo install github.com/sikhness/flogo/activity/gcpstorage
```

## Schema
Inputs and Outputs:

```json
{
  "inputs":[
    {
      "name": "jsonCredentials",
      "type": "string",
      "required": true
    },
    {
      "name": "bucketName",
      "type": "string",
      "required": true
    },
    {
      "name": "operation",
      "type": "string",
      "required": true,
      "allowed" : ["WRITE", "READ", "DELETE"]
    },
    {
      "name": "objectName",
      "type": "string",
      "required": true
    },
    {
      "name": "objectContent",
      "type": "any"
    },
    {
      "name": "writeOption",
      "type": "string",
      "value": "NEW",
      "allowed": ["NEW", "OVERWRITE", "APPEND"]
    },
    {
      "name": "objectACLList",
      "type": "params"
    }
  ],
  "outputs": [
    {
      "name": "output",
      "type": "any"
    }
  ]
}
```

## Settings
| Setting            | Required | Description |
|:---------------    |:---------|:------------|
| jsonCredentials    | True     | The service account JSON private key to access GCP |         
| bucketName         | True     | The name of the bucket for objects to be stored within GCP storage |
| operation          | True     | The operation to perform within the bucket (Allowed values are WRITE, READ, DELETE) |
| objectName         | True     | The name of the object to work with within the bucket |
| objectContent      | False    | The text content to add to the object |
| writeOption        | False    | The write option to be performed when writing an object (Allowed values are NEW, OVERWRITE, APPEND) |
| objectACLList      | False    | The ACL users and roles to assign to the object. This should follow the following JSON example: 
```json 
{
		"user1": "user-<<USER EMAIL>>",
		"role1": "OWNER",
		"user2": "user-<<USER EMAIL>>",
		"role2": "READER"
  } ```|
