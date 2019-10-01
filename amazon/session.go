package amazon

import (
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var (
	// DefaultRegion is the default region for new sessions. It defaults to "us-east-1"
	// if the environment variable "AWS_REGION" is not set.
	DefaultRegion string
)

var (
	baseSessionInstance *session.Session
	baseSessionOnce     sync.Once
)

func init() {
	DefaultRegion = os.Getenv("AWS_REGION")
	if DefaultRegion == "" {
		DefaultRegion = "us-east-1"
	}
}

// CreateSession returns a new AWS session
//
// The session is created using the AWS_REGION environment variable. If it's not
// set the default (us-east-1) region is used.
func CreateSession() (*session.Session, error) {
	return session.NewSession(&aws.Config{Region: aws.String(DefaultRegion)})
}

// MustCreateSession returns a new AWS session.
//
// If the create session call fails with an error, this function panics.
func MustCreateSession() *session.Session {
	return session.Must(CreateSession())
}

// BaseSession is the base AWS session used by all other sessions
func BaseSession() *session.Session {
	baseSessionOnce.Do(func() {
		baseSessionInstance = session.Must(CreateSession())
	})
	return baseSessionInstance
}
