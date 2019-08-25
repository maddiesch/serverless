package serverless

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Error defines an error.
type Error struct {
	Status      string                  `json:"status" validate:"required,numeric"`
	Code        string                  `json:"code" validate:"required,min=3,max=64,error_code"`
	Title       string                  `json:"title,omitempty"`
	Description string                  `json:"detail,omitempty"`
	Meta        *map[string]interface{} `json:"meta,omitempty"`
}

func (e Error) Error() string {
	data, _ := json.Marshal(e)
	return string(data)
}

// NewResponseErr creates a new error
func NewResponseErr(status int, code, description string) Error {
	return Error{
		Status:      strconv.Itoa(status),
		Code:        code,
		Title:       http.StatusText(status),
		Description: description,
	}
}

// NewResponseError returns a new ResponseError with details from the passed error
func NewResponseError(status int, code string, err error) *Error {
	new := NewResponseErr(status, code, err.Error())
	return &new
}

// NewResponseErrorFromErr creates a new response error from a generic error
func NewResponseErrorFromErr(err error) *Error {
	if re, ok := err.(Error); ok {
		return &(re)
	}
	return NewResponseError(500, "internal_server_error", err)
}

// ErrorResponse sets a JSON response pre-formatted with an error
// The first ResponseError value will be used to set the HTTP status code
func ErrorResponse(c *gin.Context, errs []*Error) error {
	if len(errs) < 1 {
		return fmt.Errorf("must provide at least one errs: got %d", len(errs))
	}

	for idx, res := range errs {
		err := GetValidator().Struct(res)
		if err != nil {
			return fmt.Errorf("err at index %d contains invalid values (%v)", idx, err)
		}
	}

	status, err := strconv.Atoi(errs[0].Status)
	if err != nil {
		return err
	}

	if status < 100 || status > 600 {
		return fmt.Errorf("err status is not a valid HTTP status code (%d)", status)
	}

	c.AbortWithStatusJSON(status, gin.H{"errors": errs})

	return nil
}

// MustErrorResponse sets a JSON response pre-formatted with an error
// The first Error value will be used to set the HTTP status code
func MustErrorResponse(c *gin.Context, errs []*Error) {
	err := ErrorResponse(c, errs)
	if err != nil {
		panic(err)
	}
}
