package urlhash_test

import (
	"testing"

	"github.com/joshturge/shortener/pkg/urlhash"
)

var testCase = map[string]string{
	"this": "45a4cd58",
	"is":   "04b90f56",
	"test": "4fdcca5d",
}

func TestHash(t *testing.T) {
	for key, val := range testCase {
		if h, err := urlhash.Hash(key); err != nil || h != val {
			t.Errorf("Hash failed: %w wanted %s got %s", err, val, h)
		}
	}
}
