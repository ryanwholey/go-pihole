package pihole

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdBlocker(t *testing.T) {
	t.Run("enable ad blocker", func(t *testing.T) {
		isAcceptance(t)

		c := newTestClient()
		ctx := context.Background()

		_, err := c.AdBlocker.Update(ctx, AdBlockerStatusOptions{
			Enabled: false,
		})
		require.NoError(t, err)

		status, err := c.AdBlocker.Update(ctx, AdBlockerStatusOptions{
			Enabled: true,
		})
		require.NoError(t, err)

		assert.Equal(t, status.Enabled, true)
	})

	t.Run("disable ad blocker", func(t *testing.T) {
		isAcceptance(t)

		c := newTestClient()
		ctx := context.Background()

		_, err := c.AdBlocker.Update(ctx, AdBlockerStatusOptions{
			Enabled: true,
		})
		require.NoError(t, err)

		status, err := c.AdBlocker.Update(ctx, AdBlockerStatusOptions{
			Enabled: false,
		})
		require.NoError(t, err)

		assert.Equal(t, status.Enabled, false)
	})
}
