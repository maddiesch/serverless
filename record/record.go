package record

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	dba "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/maddiesch/serverless/amazon"
)

// ClientProvider is the function signature for a client provider
type ClientProvider func() (*dynamodb.DynamoDB, error)

var (
	// DefaultClientProvider is the the function that will be called when a record
	// operation is needed but there is no active client
	DefaultClientProvider = func() (*dynamodb.DynamoDB, error) {
		return dynamodb.New(amazon.BaseSession()), nil
	}

	// TableName is the table name pointer for the DynamoDB table
	TableName *string
)

// Record is used to store an object into a DynamoDB table
type Record interface {
	// Should return the primary key values for the record
	Key() map[string]*dynamodb.AttributeValue
}

// ExtendedRecord allows for a Record to generate extra data during marshaling
type ExtendedRecord interface {
	Record

	Marshal(map[string]*dynamodb.AttributeValue) error

	Unmarshal(map[string]*dynamodb.AttributeValue) error
}

// Save writes the record into the DynamoDB table
func Save(ctx context.Context, r Record) error {
	item, err := Marshal(r)
	if err != nil {
		return err
	}

	_, err = Client().PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: TableName,
		Item:      item,
	})

	return err
}

// Destroy performs a destroy operation for the Record
func Destroy(ctx context.Context, r Record) error {
	key := r.Key()
	if len(key) == 0 {
		return fmt.Errorf("record has no key's")
	}

	return nil
}

// Marshal returns a DynamoDB item for the record
func Marshal(r Record) (map[string]*dynamodb.AttributeValue, error) {
	item, err := dba.MarshalMap(r)
	if err != nil {
		return nil, err
	}

	for key, value := range r.Key() {
		item[key] = value
	}

	if er, ok := r.(ExtendedRecord); ok {
		if err := er.Marshal(item); err != nil {
			return nil, err
		}
	}

	return item, nil
}

// Unmarshal converts a DynamoDB item into a Record
func Unmarshal(data map[string]*dynamodb.AttributeValue, r Record) error {
	err := dba.UnmarshalMap(data, r)
	if err != nil {
		return err
	}

	if er, ok := r.(ExtendedRecord); ok {
		if err := er.Unmarshal(data); err != nil {
			return err
		}
	}

	return nil
}

func init() {
	TableName = aws.String(os.Getenv("DYNAMODB_TABLE_NAME"))
}

var (
	clientMutex  sync.Mutex
	sharedClient *dynamodb.DynamoDB
)

// Client is the default client
func Client() *dynamodb.DynamoDB {
	clientMutex.Lock()
	defer clientMutex.Unlock()

	if sharedClient == nil {
		new, err := DefaultClientProvider()
		if err != nil {
			panic(err)
		}

		sharedClient = new
	}

	return sharedClient
}
