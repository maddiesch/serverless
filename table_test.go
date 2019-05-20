package serverless

import (
	"fmt"
	"testing"

	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/require"
)

func TestCreateDynamoTable(t *testing.T) {
	client := testDBInstance.Client
	tableName := fmt.Sprintf("testing-%s", ksuid.New().String())
	input := &CreateDynamoTableInput{
		Client:     client,
		TableName:  tableName,
		PrimaryKey: &CreateDynamoTableKey{PartitionKey: "PK", RangeKey: "SK"},
		GlobalSecondaryIndexes: map[string]*CreateDynamoTableKey{
			"GSI1": &CreateDynamoTableKey{PartitionKey: "GSI1PK", RangeKey: "GSI1SK"},
			"GSI2": &CreateDynamoTableKey{PartitionKey: "GSI2PK", RangeKey: "GSI2SK"},
		},
		LocalSecondaryIndexes: map[string]*CreateDynamoTableKey{
			"LSI1": &CreateDynamoTableKey{PartitionKey: "PK", RangeKey: "LSI1SK"},
			"LSI2": &CreateDynamoTableKey{PartitionKey: "PK", RangeKey: "LSI2SK"},
		},
	}

	err := CreateDynamoTable(input)
	require.NoError(t, err)

	err = DeleteDynamoTable(&DeleteDynamoTableInput{Client: client, TableName: tableName})
	require.NoError(t, err)
}
