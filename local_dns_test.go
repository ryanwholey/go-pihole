package pihole

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testAssertDNS(t *testing.T, c *Client, expected *DNSRecord, assertErr error) {
	actual, err := c.LocalDNS.Get(context.TODO(), expected.Domain)
	if assertErr != nil {
		assert.ErrorAs(t, err, assertErr)
		return
	}

	require.NoError(t, err)

	assert.Equal(t, expected.Domain, actual.Domain)
	assert.Equal(t, expected.IP, actual.IP)
}

func cleanupDNS(t *testing.T, c *Client, domain string) {
	if err := c.LocalDNS.Delete(context.TODO(), domain); err != nil {
		log.Printf("Failed to clean up domain record: %s\n", domain)
	}
}

func TestLocalDNS(t *testing.T) {
	t.Run("Test create a DNS record", func(t *testing.T) {
		isAcceptance(t)

		c := newTestClient(t)
		defer cleanupTestClient(c)

		domain := fmt.Sprintf("test.%s", randomID())

		record, err := c.LocalDNS.Create(context.Background(), domain, "127.0.0.1")
		require.NoError(t, err)

		defer cleanupDNS(t, c, domain)

		testAssertDNS(t, c, record, nil)
		testAssertDNS(t, c, &DNSRecord{
			Domain: record.Domain,
			IP:     "127.0.0.1",
		}, nil)
	})

	t.Run("Test delete a DNS record", func(t *testing.T) {
		isAcceptance(t)

		c := newTestClient(t)
		defer cleanupTestClient(c)

		ctx := context.Background()

		domain := fmt.Sprintf("test.%s", randomID())

		record, err := c.LocalDNS.Create(ctx, domain, "127.0.0.1")
		require.NoError(t, err)
		defer cleanupDNS(t, c, record.Domain)

		err = c.LocalDNS.Delete(ctx, domain)
		require.NoError(t, err)

		_, err = c.LocalDNS.Get(ctx, domain)
		assert.ErrorIs(t, err, ErrorLocalDNSNotFound)
	})
}
