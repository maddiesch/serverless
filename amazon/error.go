package amazon

import "github.com/aws/aws-sdk-go/aws/awserr"

// AWSError returns an AWS Error or nil
func AWSError(err error) awserr.Error {
	if err, ok := err.(awserr.Error); ok {
		return err
	}
	return nil
}

// IsErrorCode checks if an error is an AWS Error and then check if the codes match
func IsErrorCode(err error, code string) bool {
	er := AWSError(err)
	if er == nil {
		return false
	}
	return er.Code() == code
}
