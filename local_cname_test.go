package pihole

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testAssertCNAME(t *testing.T, c Client, expected *CNAMERecord, assertErr error) {
	actual, err := c.LocalCNAME.Get(context.TODO(), expected.Domain)
	if assertErr != nil {
		assert.ErrorAs(t, err, assertErr)
		return
	}

	require.NoError(t, err)

	assert.Equal(t, expected.Domain, actual.Domain)
	assert.Equal(t, expected.Target, actual.Target)
}

func cleanupCNAME(t *testing.T, c Client, domain string) {
	if err := c.LocalCNAME.Delete(context.TODO(), domain); err != nil {
		log.Printf("Failed to clean up CNAME record: %s\n", domain)
	}
}

func TestLocalCNAME(t *testing.T) {
	t.Run("Test create a CNAME record", func(t *testing.T) {
		isAcceptance(t)

		c := newTestClient()

		domain := fmt.Sprintf("test.%s", randomID())

		record, err := c.LocalCNAME.Create(context.Background(), domain, "domain.com")
		require.NoError(t, err)

		defer cleanupCNAME(t, c, domain)

		testAssertCNAME(t, c, record, nil)
		testAssertCNAME(t, c, &CNAMERecord{
			Domain: record.Domain,
			Target: "domain.com",
		}, nil)
	})

	t.Run("Test delete a CNAME record", func(t *testing.T) {
		isAcceptance(t)

		c := newTestClient()
		ctx := context.Background()

		domain := fmt.Sprintf("test.%s", randomID())

		record, err := c.LocalCNAME.Create(ctx, domain, "domain.com")
		require.NoError(t, err)
		defer cleanupCNAME(t, c, record.Domain)

		err = c.LocalCNAME.Delete(ctx, domain)
		require.NoError(t, err)

		_, err = c.LocalCNAME.Get(ctx, domain)
		assert.ErrorIs(t, err, ErrorLocalCNAMENotFound)
	})
}
