package tester

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/maddiesch/serverless"
	"github.com/maddiesch/serverless/record"
)

var (
	// EnvironmentFilePath is a path to the environment
	EnvironmentFilePath string
)

// LoadEnvironment will load an environment JSON file from an AWS-SAM env JSON
func LoadEnvironment(path, key string) error {
	_, loc, _, ok := runtime.Caller(1)
	if !ok {
		return errors.New("Failed to load caller path")
	}

	body, err := ioutil.ReadFile(filepath.Join(filepath.Dir(loc), path))
	if err != nil {
		return err
	}

	content := map[string]map[string]string{}
	err = json.Unmarshal(body, &content)
	if err != nil {
		return err
	}

	env, ok := content[key]
	if !ok {
		return fmt.Errorf("failed to find a environment for %s", key)
	}

	for key, value := range env {
		if os.Getenv(key) != "" {
			continue
		}
		os.Setenv(key, value)
	}

	return nil
}

// Run performs any setup necessary for the test suite
func Run(m *testing.M) int {
	os.Setenv("AWS_SAM_LOCAL", "true")
	os.Setenv("IS_TESTING", "true")

	record.TableName = aws.String(os.Getenv("DYNAMODB_TABLE_NAME"))

	return m.Run()
}

// Request performs a request into he app
func Request(app serverless.Application, r *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	err := app.Router().Dispatch(context.Background(), w, r)
	if err != nil {
		panic(err)
	}
	return w
}
