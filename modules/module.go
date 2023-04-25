package modules

import "sync"

// Used by external Sources
func NewSourceBase(options *Options) *ModuleBase {
	return &ModuleBase{
		Flatmap: make(map[string]any),
		RWTex:   sync.RWMutex{},
		Options: options,
	}
}

// Used by Configuration
func (module *ModuleBase) Connect(refreshC chan Module) {
	module.RefreshC = refreshC
}

// Checks if a key exists
func (module *ModuleBase) Exists(key string) bool {
	return module.Get(key) != nil
}

// Get Config Values
func (module *ModuleBase) Get(key string) (value interface{}) {
	if value, ok := module.Flatmap[key]; ok {
		return value
	}

	return nil
}

func (module *ModuleBase) GetKeys() (result []string) {
	result = make([]string, 0)
	for key := range module.Flatmap {
		result = append(result, key)
	}

	return result
}

func (module *ModuleBase) GetOptions() *Options {
	return module.Options
}

func (module *ModuleBase) NotifyDirtyness(externalSource Module) {
	if module.RefreshC != nil {
		module.RefreshC <- externalSource
	}
}
