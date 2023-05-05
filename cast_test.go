package configurationmanager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCastAndTryAssignValue(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "cast string to integer should return integer",
			run: func(t *testing.T) {
				from := "1"
				var to int
				CastAndAssignValue(from, &to)
				assert.Equal(t, 1, to)
			},
		},
		{
			name: "cast string to integer pointer should return integer pointer",
			run: func(t *testing.T) {
				from := pointer("1")
				var to *int
				CastAndAssignValue(&from, &to)
				assert.Equal(t, 1, *to)
			},
		},
		{
			name: "cast string to nested integer pointer should return integer pointer",
			run: func(t *testing.T) {
				from := pointer("1")
				var to ******int
				CastAndAssignValue(&from, &to)
				assert.Equal(t, 1, ******to)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, test.run)
	}
}
