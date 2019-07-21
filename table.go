package serverless

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// CreateDynamoTableInput contains the required values for a create table call
type CreateDynamoTableInput struct {
	Client                 *dynamodb.DynamoDB    `validate:"required"`
	TableName              string                `validate:"min=1"`
	PrimaryKey             *CreateDynamoTableKey `validate:"required"`
	GlobalSecondaryIndexes map[string]*CreateDynamoTableKey
	LocalSecondaryIndexes  map[string]*CreateDynamoTableKey
}

// CreateDynamoTableKey contains the partition & range key pair
type CreateDynamoTableKey struct {
	PartitionKey string `validate:"min=2"`
	RangeKey     string `validate:"min=2"`
}

func (k *CreateDynamoTableKey) attributes() []*dynamodb.AttributeDefinition {
	return []*dynamodb.AttributeDefinition{
		&dynamodb.AttributeDefinition{
			AttributeName: aws.String(k.PartitionKey),
			AttributeType: aws.String("S"),
		},
		&dynamodb.AttributeDefinition{
			AttributeName: aws.String(k.RangeKey),
			AttributeType: aws.String("S"),
		},
	}
}

func (k *CreateDynamoTableKey) element() []*dynamodb.KeySchemaElement {
	return []*dynamodb.KeySchemaElement{
		&dynamodb.KeySchemaElement{
			AttributeName: aws.String(k.PartitionKey),
			KeyType:       aws.String("HASH"),
		},
		&dynamodb.KeySchemaElement{
			AttributeName: aws.String(k.RangeKey),
			KeyType:       aws.String("RANGE"),
		},
	}
}

// CreateDynamoTable creates a table and waits for the complete to finish
func CreateDynamoTable(input *CreateDynamoTableInput) error {
	if err := GetValidator().Struct(input); err != nil {
		return err
	}

	params := &dynamodb.CreateTableInput{
		TableName:            aws.String(input.TableName),
		KeySchema:            make([]*dynamodb.KeySchemaElement, 0),
		AttributeDefinitions: make([]*dynamodb.AttributeDefinition, 0),
		BillingMode:          aws.String("PROVISIONED"),
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	}

	params.AttributeDefinitions = append(params.AttributeDefinitions, input.PrimaryKey.attributes()...)
	params.KeySchema = input.PrimaryKey.element()

	if len(input.GlobalSecondaryIndexes) > 0 {
		for name, idx := range input.GlobalSecondaryIndexes {
			params.AttributeDefinitions = append(params.AttributeDefinitions, idx.attributes()...)
			params.GlobalSecondaryIndexes = append(params.GlobalSecondaryIndexes, &dynamodb.GlobalSecondaryIndex{
				IndexName:  aws.String(name),
				KeySchema:  idx.element(),
				Projection: &dynamodb.Projection{ProjectionType: aws.String("ALL")},
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(5),
					WriteCapacityUnits: aws.Int64(5),
				},
			})
		}
	}

	if len(input.LocalSecondaryIndexes) > 0 {
		for name, idx := range input.LocalSecondaryIndexes {
			params.AttributeDefinitions = append(params.AttributeDefinitions, idx.attributes()[1])
			params.LocalSecondaryIndexes = append(params.LocalSecondaryIndexes, &dynamodb.LocalSecondaryIndex{
				IndexName:  aws.String(name),
				KeySchema:  idx.element(),
				Projection: &dynamodb.Projection{ProjectionType: aws.String("KEYS_ONLY")},
			})
		}
	}

	_, err := input.Client.CreateTable(params)
	if err != nil {
		return err
	}

	return input.Client.WaitUntilTableExists(&dynamodb.DescribeTableInput{
		TableName: aws.String(input.TableName),
	})
}

// DeleteDynamoTableInput contains the values needed for a DeleteDynamoTable call
type DeleteDynamoTableInput struct {
	Client    *dynamodb.DynamoDB
	TableName string
}

// DeleteDynamoTable performs a delete of the table and waits for it to be completely deleted.
func DeleteDynamoTable(input *DeleteDynamoTableInput) error {
	params := &dynamodb.DeleteTableInput{
		TableName: aws.String(input.TableName),
	}

	_, err := input.Client.DeleteTable(params)
	if err != nil {
		return err
	}

	return input.Client.WaitUntilTableNotExists(&dynamodb.DescribeTableInput{
		TableName: aws.String(input.TableName),
	})
}
