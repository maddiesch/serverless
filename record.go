package serverless

import "github.com/aws/aws-sdk-go/service/dynamodb"

// Record is an interface for managing DynamoDB records.
type Record interface {
	Attributes() (map[string]*dynamodb.AttributeValue, error)

	Key() (map[string]*dynamodb.AttributeValue, error)

	Assign(map[string]*dynamodb.AttributeValue) error
}
