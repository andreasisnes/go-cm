package configurationmanager

import (
	"testing"

	"github.com/andreasisnes/go-configuration-manager/modules"
	"github.com/stretchr/testify/assert"
)

func TestBuilderNew(t *testing.T) {
	RunTests(t, &[]Test{
		{
			name: "should return non nil value",
			run: func(t *testing.T) {
				builder := New(nil)
				assert.NotNil(t, builder)
				assert.NotNil(t, builder.Modules())
			},
		},
	})
}

func TestBuilderAdd(t *testing.T) {
	RunTests(t, &[]Test{
		{
			name: "given nil argument should not be added",
			run: func(t *testing.T) {
				builder := builder{
					modules: []modules.Module{},
				}
				builder.Add(nil)
				assert.Empty(t, builder.modules)
			},
		},
		{
			name: "given non nil argument should be added to list",
			run: func(t *testing.T) {
				builder := builder{
					modules: []modules.Module{},
				}

				builder.Add(&testModule{})
				builder.Add(&testModule{})
				assert.Len(t, builder.modules, 2)
			},
		},
	})
}

func TestBuilderClear(t *testing.T) {
	RunTests(t, &[]Test{
		{
			name: "should create a new list",
			run: func(t *testing.T) {
				builder := builder{
					modules: []modules.Module{},
				}

				builder.Add(&testModule{})
				builder.Clear()
				assert.Empty(t, builder.modules)
			},
		},
	})
}

func TestBuilderModules(t *testing.T) {
	RunTests(t, &[]Test{
		{
			name: "should return list of modules",
			run: func(t *testing.T) {
				builder := builder{
					modules: []modules.Module{
						&testModule{},
					},
				}
				assert.Len(t, builder.Modules(), 1)
			},
		},
	})
}

func TestBuilderBuild(t *testing.T) {
	RunTests(t, &[]Test{
		{
			name: "should return a configuration instance",
			run: func(t *testing.T) {
				config := New(nil).
					Build()

				assert.NotNil(t, config)
			},
		},
	})
}
