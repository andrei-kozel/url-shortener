package shorten

import (
	"testing"

	"github.com/andrei-kozel/url-shortener/internal/shorten"
	"github.com/stretchr/testify/assert"
)

func TestShorten(t *testing.T) {
	t.Run("Returns an alphanumeric short identifier", func(t *testing.T) {
		type testCase struct {
			id       uint32
			expected string
		}

		testCases := []testCase{
			{id: 1024, expected: "Mv"},
			{id: 0, expected: ""},
			{id: 1, expected: "y"},
		}

		for _, tc := range testCases {
			actual := shorten.Shorten(tc.id)
			assert.Equal(t, tc.expected, actual)
		}
	})
	t.Run("Is identical for the same input", func(t *testing.T) {})
}
