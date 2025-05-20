package pihole

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientValidation(t *testing.T) {

	t.Run("no error on valid client config", func(t *testing.T) {
		isUnit(t)
		t.Parallel()

		_, err := New(Config{
			BaseURL:  "http://localhost:8080",
			Password: "test",
		})

		assert.NoError(t, err)
	})
}

func isAcceptance(t *testing.T) {
	if os.Getenv("TEST_ACC") != "1" {
		t.Skip("skipping acceptance test")
	} else {
		accPreflghtCheck(t)
	}
}

func isUnit(t *testing.T) {
	if os.Getenv("TEST_ACC") == "1" {
		t.Skip("skipping unit test")
	}
}

func accPreflghtCheck(t *testing.T) {
	require.NotEmpty(t, os.Getenv("PIHOLE_URL"), "PIHOLE_URL must be set for acceptance tests")
	require.NotEmpty(t, os.Getenv("PIHOLE_PASSWORD"), "PIHOLE_PASSWORD must be set for acceptance tests")
}

func newTestClient(t *testing.T) *Client {
	c, err := New(Config{
		BaseURL:  os.Getenv("PIHOLE_URL"),
		Password: os.Getenv("PIHOLE_PASSWORD"),
	})

	require.NoError(t, err)

	return c
}

func randomID() string {
	b := make([]byte, 5)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Errorf("failed to make random ID: %w", err))
	}

	return fmt.Sprintf("%X", b)
}

func cleanupTestClient(c *Client) {
	if err := c.SessionAPI.Logout(context.TODO()); err != nil {
		fmt.Printf("failed to clean up client after acceptance test: %s\n", err)
	}
}
