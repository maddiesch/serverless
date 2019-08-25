package serverless

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorCodeValidation(t *testing.T) {
	t.Run("given a valid key", func(t *testing.T) {
		value := struct {
			Code string `validate:"error_code"`
		}{Code: "foo_Bar_BAZ_1"}

		err := GetValidator().Struct(value)

		require.NoError(t, err)
	})

	t.Run("given an invalid key with a leading underscore", func(t *testing.T) {
		value := struct {
			Code string `validate:"error_code"`
		}{Code: "_foo_Bar_BAZ"}

		err := GetValidator().Struct(value)

		require.Error(t, err)
	})

	t.Run("given an invalid key with a dash", func(t *testing.T) {
		value := struct {
			Code string `validate:"error_code"`
		}{Code: "foo-Bar_BAZ"}

		err := GetValidator().Struct(value)

		require.Error(t, err)
	})

	t.Run("given an invalid key with a leading number", func(t *testing.T) {
		value := struct {
			Code string `validate:"error_code"`
		}{Code: "1foo_Bar_BAZ"}

		err := GetValidator().Struct(value)

		require.Error(t, err)
	})
}
