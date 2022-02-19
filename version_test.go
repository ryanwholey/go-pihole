package pihole

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersion(t *testing.T) {
	t.Run("Fetch versions", func(t *testing.T) {
		isAcceptance(t)

		c := newTestClient()

		versions, err := c.Version.Get(context.Background())
		require.NoError(t, err)

		assert.NotEmpty(t, versions.CoreCurrent)
		assert.NotEmpty(t, versions.FTLCurrent)
		assert.NotEmpty(t, versions.WebCurrent)
	})
}
