package pihole

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuth(t *testing.T) {
	tcs := []struct {
		name string
	}{
		{
			name: "Sets SID and CSRF on the client",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			c := newTestClient(t)

			require.NoError(t, c.AuthAPI.Authenticate(context.TODO()))

			assert.NotEmpty(t, c.auth.csrf)
			assert.NotEmpty(t, c.auth.sid)
		})
	}
}
