package pihole

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCustomDNS(t *testing.T) {
	t.Run("list custom DNS records", func(t *testing.T) {
		client := New(Config{
			BaseURL:  os.Getenv("PIHOLE_URL"),
			APIToken: os.Getenv("PIHOLE_API_TOKEN"),
		})

		list, err := client.CustomDNS.List(context.Background())
		require.NoError(t, err)

		assert.ElementsMatch(t, list, CustomDNSList{
			{
				Domain: "ryan.com",
				IP:     "127.0.0.1",
			},
		})
	})

	t.Run("create custom DNS records", func(t *testing.T) {
		client := New(Config{
			BaseURL:  os.Getenv("PIHOLE_URL"),
			APIToken: os.Getenv("PIHOLE_API_TOKEN"),
		})

		record, err := client.CustomDNS.Create(context.Background(), "ryanwholey.com", "127.0.0.1")
		require.NoError(t, err)
		assert.Equal(t, CustomDNS{Domain: "ryanwholey.com", IP: "127.0.0.1"}, record)
	})
}
