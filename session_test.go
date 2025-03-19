package pihole

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionLogin(t *testing.T) {
	tcs := []struct {
		name string
	}{
		{
			name: "Login SID on the client",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			isAcceptance(t)

			c := newTestClient(t)

			_, err := c.SessionAPI.Login(context.TODO())
			require.NoError(t, err)

			assert.NotEmpty(t, c.auth.sid)
		})
	}
}
