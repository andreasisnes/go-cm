package configurationmanager

import (
	"testing"

	"github.com/andreasisnes/go-configuration-manager/modules"
)

type Test struct {
	name string
	run  func(t *testing.T)
}

func RunTests(t *testing.T, tests *[]Test) {
	for _, test := range *tests {
		t.Run(test.name, test.run)
	}
}

func newTestModule(options *modules.Options, opts ...func(map[string]any)) *testModule {
	if options == nil {
		options = &modules.Options{}
	}
	t := &testModule{
		values:     make(map[string]any),
		ModuleBase: modules.NewSourceBase(options),
	}

	for _, opt := range opts {
		opt(t.values)
		opt(t.Flatmap)
	}

	return t
}

type testModule struct {
	*modules.ModuleBase
	values map[string]any
}

func (module *testModule) GetRefreshedValue(key string) any {
	return nil
}

func (module *testModule) Deconstruct() {

}

func (module *testModule) Load() {

}

func (module *testModule) Exists(key string) bool {
	_, ok := module.values[key]
	return ok
}

func (module *testModule) Get(key string) any {
	if val, ok := module.values[key]; ok {
		return val
	}

	return nil
}
