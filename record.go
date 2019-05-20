package serverless

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Record is an interface for managing DynamoDB records.
type Record interface {
	Attributes() (map[string]*dynamodb.AttributeValue, error)

	Key() (map[string]*dynamodb.AttributeValue, error)

	Assign(map[string]*dynamodb.AttributeValue) error
}

// MarshalRecord turns a record into a DynamoDB record
func MarshalRecord(r Record) (map[string]*dynamodb.AttributeValue, error) {
	attr, err := r.Attributes()
	if err != nil {
		return map[string]*dynamodb.AttributeValue{}, err
	}

	key, err := r.Key()
	if err != nil {
		return map[string]*dynamodb.AttributeValue{}, err
	}

	for key, value := range key {
		attr[key] = value
	}

	return attr, nil
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

	item, err := MarshalRecord(r)
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

	item, err := MarshalRecord(r)
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

	item, err := MarshalRecord(r)
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
