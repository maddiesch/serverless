package amazon

import (
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

const (
	// DefaultRegion is the default region for AWS configuration.
	DefaultRegion = "us-east-1"
)

var (
	baseSessionInstance *session.Session
	baseSessionOnce     sync.Once
)

// CreateSession returns a new AWS session
//
// The session is created using the AWS_REGION environment variable. If it's not
// set the default (us-east-1) region is used.
func CreateSession() (*session.Session, error) {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}
	return session.NewSession(&aws.Config{Region: aws.String(region)})
}

// MustCreateSession returns a new AWS session.
//
// If the create session call fails with an error, this function panics.
func MustCreateSession() *session.Session {
	ses, err := CreateSession()
	if err != nil {
		panic(err)
	}
	return ses
}

// BaseSession is the base AWS session used by all other sessions
func BaseSession() *session.Session {
	baseSessionOnce.Do(func() {
		baseSessionInstance = MustCreateSession()
	})
	return baseSessionInstance
}
