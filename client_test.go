package pihole

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientValidation(t *testing.T) {
	t.Run("error on unset API token", func(t *testing.T) {
		t.Parallel()
		c := New(Config{
			BaseURL: "http://localhost:8080",
		})

		assert.ErrorIs(t, c.Validate(), ErrClientValidation)
	})

	t.Run("error on unset URL", func(t *testing.T) {
		t.Parallel()
		c := New(Config{
			APIToken: "token",
		})

		assert.ErrorIs(t, c.Validate(), ErrClientValidation)
	})

	t.Run("no error on valid client config", func(t *testing.T) {
		t.Parallel()
		c := New(Config{
			BaseURL:  "http://localhost:8080",
			APIToken: "token",
		})

		assert.NoError(t, c.Validate())
	})
}
