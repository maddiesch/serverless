package serverless

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/maddiesch/serverless/amazon"
	"github.com/maddiesch/serverless/sam"
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
		ses.Config.Endpoint = aws.String(DefaultEnv("AWS_DYNAMODB_ENDPOINT", "http://docker.for.mac.localhost:8000"))
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

var (
	// ErrRecordAlreadyExist is returned when a `CreateRecord` call fails because the item exists.
	ErrRecordAlreadyExist = errors.New("record failed to create because it already exists")

	// ErrRecordDoesNotExist is returned when a `UpdateRecord` call fails because the item doesn't exist.
	ErrRecordDoesNotExist = errors.New("record failed to update because it does not exist")

	// ErrRecordNotFound is returned when a single record can not be found.
	ErrRecordNotFound = errors.New("record does not exist")
)

// CreateRecord writes the record to the table if it doesn't already exist.
func (db *DB) CreateRecord(r Record) error {
	valid, validate := r.(Validatable)
	if validate {
		err := valid.Validate()
		if err != nil {
			return err
		}
	}

	item, err := createFullItem(r)
	if err != nil {
		return err
	}

	ce := fmt.Sprintf("attribute_not_exists(%s)", db.createKeyName)

	params := &dynamodb.PutItemInput{
		TableName:           db.tn(),
		ConditionExpression: aws.String(ce),
		Item:                item,
	}

	_, err = db.Client.PutItem(params)

	if IsConditionalCheckFailure(err) {
		return ErrRecordAlreadyExist
	}

	return err
}

// UpdateRecord updates the record in the table, only if it already exists.
func (db *DB) UpdateRecord(r Record) error {
	valid, validate := r.(Validatable)
	if validate {
		err := valid.Validate()
		if err != nil {
			return err
		}
	}

	item, err := createFullItem(r)
	if err != nil {
		return err
	}

	ce := fmt.Sprintf("attribute_exists(%s)", db.createKeyName)

	params := &dynamodb.PutItemInput{
		TableName:           db.tn(),
		ConditionExpression: aws.String(ce),
		Item:                item,
	}

	_, err = db.Client.PutItem(params)

	if IsConditionalCheckFailure(err) {
		return ErrRecordDoesNotExist
	}

	return err
}

// SaveRecord creates or updates the record in the table.
func (db *DB) SaveRecord(r Record) error {
	valid, validate := r.(Validatable)
	if validate {
		err := valid.Validate()
		if err != nil {
			return err
		}
	}

	item, err := createFullItem(r)
	if err != nil {
		return err
	}

	params := &dynamodb.PutItemInput{
		TableName: db.tn(),
		Item:      item,
	}

	_, err = db.Client.PutItem(params)

	return err
}

// DestroyRecord deletes the item
func (db *DB) DestroyRecord(r Record) error {
	valid, validate := r.(Validatable)
	if validate {
		err := valid.Validate()
		if err != nil {
			return err
		}
	}

	key, err := r.Key()
	if err != nil {
		return err
	}

	input := &dynamodb.DeleteItemInput{
		TableName: db.tn(),
		Key:       key,
	}

	_, err = db.Client.DeleteItem(input)

	return err
}

// IsConditionalCheckFailure returns true if the error passed was caused by a ConditionalCheckFailedException
func IsConditionalCheckFailure(err error) bool {
	if ae, ok := err.(awserr.RequestFailure); ok && ae.Code() == "ConditionalCheckFailedException" {
		return true
	}
	return false
}

func createFullItem(r Record) (map[string]*dynamodb.AttributeValue, error) {
	key, err := r.Key()
	if err != nil {
		return map[string]*dynamodb.AttributeValue{}, err
	}

	attr, err := r.Attributes()
	if err != nil {
		return map[string]*dynamodb.AttributeValue{}, err
	}

	for k, v := range key {
		attr[k] = v
	}

	return attr, nil
}
