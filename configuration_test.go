package configurationmanager

import (
	"testing"

	"github.com/andreasisnes/go-configuration-manager/modules"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	RunTests(t, &[]Test{
		{
			name: "given two layer with same variable, should return top layer value",
			run: func(t *testing.T) {
				config := New(nil).
					Add(newTestModule(nil, func(m map[string]any) {
						m["test-int"] = 10
					})).
					Add(newTestModule(nil, func(m map[string]any) {
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
					Add(newTestModule(nil)).
					Build()

				var strPtr *string = nil
				val := config.Get("test", &strPtr)
				assert.Nil(t, val)
				assert.Nil(t, strPtr)
			},
		},
	})
}

func TestList(t *testing.T) {
	RunTests(t, &[]Test{
		{
			name: "",
			run: func(t *testing.T) {
				config := New(&Options{
					Delimiter: ":",
				}).
					Add(newTestModule(&modules.Options{
						Delimiter: "-",
					}, func(m map[string]any) {
						m["test-int1"] = 10
					})).
					Add(newTestModule(&modules.Options{
						Delimiter: "?",
					}, func(m map[string]any) {
						m["test?int2"] = 100
					})).
					Build()

				result := config.List()
				assert.Contains(t, result, "test:int1")
				assert.Contains(t, result, "test:int2")
			},
		},
	})
}
