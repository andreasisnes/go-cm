package configurationmanager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	RunTests(t, &[]Test{
		{
			name: "given two layer with same variable, should return top layer value",
			run: func(t *testing.T) {
				config := New(nil).
					Add(newTestModule(func(m map[string]any) {
						m["test-int"] = 10
					})).
					Add(newTestModule(func(m map[string]any) {
						m["test-int"] = 100
					})).
					Build()

				var strPtr *string
				config.Get("test-int", &strPtr)
				assert.Equal(t, "100", *strPtr)
			},
		},
		{
			name: "given a layer with a variable, should return nil if not found",
			run: func(t *testing.T) {
				config := New(nil).
					Add(newTestModule()).
					Build()

				var strPtr *string = nil
				val := config.Get("test", &strPtr)
				assert.Nil(t, val)
				assert.Nil(t, strPtr)
			},
		},
	})
}

func TestUnmarshal(t *testing.T) {
	RunTests(t, &[]Test{
		{
			name: "given flat struct with int and string values should umarshal",
			run: func(t *testing.T) {
				type Unmarshal struct {
					Test int
					Name string
				}

				config := New(nil).
					Add(newTestModule(func(m map[string]any) {
						m["Name"] = "test"
						m["Test"] = "100"
					})).Build()

				result := Unmarshal{}
				config.Unmarshal(&result)
				assert.Equal(t, "test", result.Name)
				assert.Equal(t, 100, result.Test)
			},
		},
	})
}
