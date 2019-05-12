package serverless

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/require"
)

type testRecord struct {
	PK    string
	SK    string
	Value string
}

func (r *testRecord) Attributes() (map[string]*dynamodb.AttributeValue, error) {
	return map[string]*dynamodb.AttributeValue{
		"Value": {S: aws.String(r.Value)},
	}, nil
}

func (r *testRecord) Key() (map[string]*dynamodb.AttributeValue, error) {
	return map[string]*dynamodb.AttributeValue{
		"PK": {S: aws.String(r.PK)},
		"SK": {S: aws.String(r.SK)},
	}, nil
}

func (r *testRecord) Assign(v map[string]*dynamodb.AttributeValue) error {
	return nil
}

func TestCreateRecord(t *testing.T) {
	record := &testRecord{PK: "primary_key", SK: ksuid.New().String(), Value: "My Test Record"}

	err := testDB().CreateRecord(record)
	require.NoError(t, err)

	err = testDB().UpdateRecord(record)
	require.NoError(t, err)
}
