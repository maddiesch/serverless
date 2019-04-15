package serverless

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorCodeValidation_valid(t *testing.T) {
	value := struct {
		Code string `validate:"error_code"`
	}{Code: "foo_Bar_BAZ_1"}

	err := GetValidator().Struct(value)

	require.NoError(t, err)
}

func TestErrorCodeValidation_invalid1(t *testing.T) {
	value := struct {
		Code string `validate:"error_code"`
	}{Code: "_foo_Bar_BAZ"}

	err := GetValidator().Struct(value)

	require.Error(t, err)
}

func TestErrorCodeValidation_invalid2(t *testing.T) {
	value := struct {
		Code string `validate:"error_code"`
	}{Code: "foo-Bar_BAZ"}

	err := GetValidator().Struct(value)

	require.Error(t, err)
}

func TestErrorCodeValidation_invalid3(t *testing.T) {
	value := struct {
		Code string `validate:"error_code"`
	}{Code: "1foo_Bar_BAZ"}

	err := GetValidator().Struct(value)

	require.Error(t, err)
}
