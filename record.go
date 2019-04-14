package serverless

import "github.com/aws/aws-sdk-go/service/dynamodb"

// Record is an interface for managing DynamoDB records.
type Record interface {
	Attributes() map[string]*dynamodb.AttributeValue

	Key() map[string]*dynamodb.AttributeValue

	Assign(map[string]*dynamodb.AttributeValue) error
}
