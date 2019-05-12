package sam

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsLocal(t *testing.T) {
	err := os.Setenv("AWS_SAM_LOCAL", "true")

	require.NoError(t, err)

	assert.Equal(t, true, IsLocal())
}
