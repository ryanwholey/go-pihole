package pihole

import (
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientValidation(t *testing.T) {
	t.Run("error on unset API token", func(t *testing.T) {
		isUnit(t)
		t.Parallel()

		c := New(Config{
			BaseURL: "http://localhost:8080",
		})

		assert.ErrorIs(t, c.Validate(), ErrClientValidation)
	})

	t.Run("error on unset URL", func(t *testing.T) {
		isUnit(t)
		t.Parallel()

		c := New(Config{
			APIToken: "token",
		})

		assert.ErrorIs(t, c.Validate(), ErrClientValidation)
	})

	t.Run("no error on valid client config", func(t *testing.T) {
		isUnit(t)
		t.Parallel()

		c := New(Config{
			BaseURL:  "http://localhost:8080",
			APIToken: "token",
		})

		assert.NoError(t, c.Validate())
	})
}

func isAcceptance(t *testing.T) {
	if os.Getenv("TEST_ACC") != "1" {
		t.Skip("skipping acceptance test")
	}
}

func isUnit(t *testing.T) {
	if os.Getenv("TEST_ACC") == "1" {
		t.Skip("skipping unit test")
	}
}

func accPreflghtCheck(t *testing.T) {
	log.Println("os.Getenv(\"PIHOLE_API_TOKEN\")")
	log.Println(os.Getenv("PIHOLE_API_TOKEN"))

	require.NotEmpty(t, os.Getenv("PIHOLE_URL"), "PIHOLE_URL must be set for acceptance tests")
	require.NotEmpty(t, os.Getenv("PIHOLE_API_TOKEN"), "PIHOLE_API_TOKEN must be set for acceptance tests")
}

func newTestClient() Client {
	return *New(Config{
		BaseURL:  os.Getenv("PIHOLE_URL"),
		APIToken: os.Getenv("PIHOLE_API_TOKEN"),
	})
}

func randomID() string {
	b := make([]byte, 5)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Errorf("failed to make random ID: %w", err))
	}

	return fmt.Sprintf("%X", b)
}
