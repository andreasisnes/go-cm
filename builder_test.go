package configurationmanager

import (
	"testing"

	"github.com/andreasisnes/go-configuration-manager/modules"
	"github.com/stretchr/testify/assert"
)

type testModule struct {
	*modules.ModuleBase
}

func (module *testModule) GetRefreshedValue(key string) any {
	return nil
}

func (module *testModule) Deconstruct() {

}

func (module *testModule) Load() {

}

func TestBuilderNew(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "should return non nil value",
			run: func(t *testing.T) {
				builder := New()
				assert.NotNil(t, builder)
				assert.NotNil(t, builder.Modules())
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, test.run)
	}
}

func TestBuilderAdd(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
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
	}

	for _, test := range tests {
		t.Run(test.name, test.run)
	}
}

func TestBuilderClear(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
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
	}

	for _, test := range tests {
		t.Run(test.name, test.run)
	}
}

func TestBuilderModules(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
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
	}

	for _, test := range tests {
		t.Run(test.name, test.run)
	}
}

func TestBuilderBuild(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "should return a configuration instance",
			run: func(t *testing.T) {
				config := New().
					Build()

				assert.NotNil(t, config)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, test.run)
	}
}
