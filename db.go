package serverless

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/maddiesch/serverless/amazon"
	"github.com/maddiesch/serverless/sam"
)

const (
	// ServerlessDynamoDBEndpointEnv is the env variable that stores the DynamoDB endpoint
	ServerlessDynamoDBEndpointEnv = "AWS_DYNAMODB_ENDPOINT"

	// ServerlessDynamoDBEndpointDefault is the default value for the endpoint if we're running locally
	ServerlessDynamoDBEndpointDefault = "http://docker.for.mac.localhost:8000"
)

// DB contains a DynamoDB config
type DB struct {
	Client        *dynamodb.DynamoDB
	TableName     string
	createKeyName string
}

// NewDB returns a new DB instance with a default session.
func NewDB(tn string, kn string) *DB {
	ses := amazon.BaseSession().Copy()

	if sam.IsLocal() {
		ses.Config.Endpoint = aws.String(GetenvDefault(ServerlessDynamoDBEndpointEnv, ServerlessDynamoDBEndpointDefault))
	}

	return NewDBWithSession(ses, tn, kn)
}

// NewDBWithSession returns a new DB instance configured with the passed in session.
func NewDBWithSession(ses *session.Session, tn string, kn string) *DB {
	db := dynamodb.New(ses)

	return &DB{
		Client:        db,
		TableName:     tn,
		createKeyName: kn,
	}
}

func (db *DB) tn() *string {
	return aws.String(db.TableName)
}

// IsConditionalCheckFailure returns true if the error passed was caused by a ConditionalCheckFailedException
func IsConditionalCheckFailure(err error) bool {
	if ae, ok := err.(awserr.RequestFailure); ok && ae.Code() == "ConditionalCheckFailedException" {
		return true
	}
	return false
}
