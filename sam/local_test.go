package sam

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsLocal(t *testing.T) {
	t.Run("given a true env", func(t *testing.T) {
		err := os.Setenv("AWS_SAM_LOCAL", "true")

		require.NoError(t, err)

		assert.Equal(t, true, IsLocal())
	})

	t.Run("given a false env", func(t *testing.T) {
		err := os.Setenv("AWS_SAM_LOCAL", "false")

		require.NoError(t, err)

		assert.Equal(t, false, IsLocal())
	})

	t.Run("given a nil env", func(t *testing.T) {
		err := os.Unsetenv("AWS_SAM_LOCAL")

		require.NoError(t, err)

		assert.Equal(t, false, IsLocal())
	})
}
