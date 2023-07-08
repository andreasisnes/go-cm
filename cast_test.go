package configurationmanager

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCastAndTryAssignValue(t *testing.T) {
	RunTests(t, &[]Test{
		{
			name: "cast string to interface{} (string) should return string",
			run: func(t *testing.T) {
				from := "test"
				var to interface{} = ""
				CastAndAssignValue(from, &to)

				assert.Equal(t, "test", to)
			},
		},
		{
			name: "cast string to time.duration should return time.duration",
			run: func(t *testing.T) {
				from := "1"
				var to time.Duration
				CastAndAssignValue(from, &to)
				assert.Equal(t, time.Nanosecond, to)
			},
		},
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
				from := "1"
				var to ******int
				CastAndAssignValue(&from, &to)
				assert.Equal(t, 1, ******to)
			},
		},
		{
			name: "cast string to timestamp should return timestamp",
			run: func(t *testing.T) {
				from := time.Now()
				var to time.Time
				CastAndAssignValue(from.Format(time.RFC3339Nano), &to)

				assert.True(t, time.Duration(math.Abs(float64(to.Sub(from)))) < time.Millisecond)
			},
		},
	})
}
