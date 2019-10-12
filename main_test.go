package serverless

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/segmentio/ksuid"

	"github.com/maddiesch/serverless/logger"
)

var (
	testDBNonce    sync.Once
	testDBInstance *DB
)

func testDB() *DB {
	testDBNonce.Do(func() {
		os.Setenv("AWS_SAM_LOCAL", "true")

		testDBInstance = NewDB(fmt.Sprintf("serverless-test-table-%s", ksuid.New().String()), "SK")

		logger.Print("Create Test Table: ", testDBInstance.TableName)

		table := &dynamodb.CreateTableInput{
			TableName: testDBInstance.tn(),
			AttributeDefinitions: []*dynamodb.AttributeDefinition{
				&dynamodb.AttributeDefinition{AttributeName: aws.String("PK"), AttributeType: aws.String("S")},
				&dynamodb.AttributeDefinition{AttributeName: aws.String("SK"), AttributeType: aws.String("S")},
			},
			KeySchema: []*dynamodb.KeySchemaElement{
				&dynamodb.KeySchemaElement{AttributeName: aws.String("PK"), KeyType: aws.String("HASH")},
				&dynamodb.KeySchemaElement{AttributeName: aws.String("SK"), KeyType: aws.String("SORT")},
			},
			ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
				WriteCapacityUnits: aws.Int64(2),
				ReadCapacityUnits:  aws.Int64(2),
			},
		}

		_, err := testDBInstance.Client.CreateTable(table)
		if err != nil {
			panic(err)
		}

		testDBInstance.Client.WaitUntilTableExists(&dynamodb.DescribeTableInput{TableName: testDBInstance.tn()})
	})

	return testDBInstance
}

func teardownDB() {
	logger.Print("Delete Test Table: ", testDBInstance.TableName)

	_, err := testDB().Client.DeleteTable(&dynamodb.DeleteTableInput{TableName: testDB().tn()})
	if err != nil {
		panic(err)
	}

	testDBInstance.Client.WaitUntilTableNotExists(&dynamodb.DescribeTableInput{TableName: testDBInstance.tn()})
}

func TestMain(m *testing.M) {
	testDB() // Setup DB
	status := m.Run()
	teardownDB()
	os.Exit(status)
}
