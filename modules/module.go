package modules

import "sync"

// Used by external Sources
func NewSourceBase(options *Options) *ModuleBase {
	return &ModuleBase{
		Flatmap: make(map[string]interface{}),
		RWTex:   sync.RWMutex{},
		Options: options,
	}
}

// Used by Configuration
func (module *ModuleBase) Connect(refreshC chan Module) {
	module.RWTex.Lock()
	defer module.RWTex.Unlock()
	module.RefreshC = refreshC
}

// Checks if a key exists
func (module *ModuleBase) Exists(key string) bool {
	return module.Get(key) != nil
}

// Get Config Values
func (module *ModuleBase) Get(key string) (value interface{}) {
	module.RWTex.RLock()
	defer module.RWTex.RUnlock()

	if value, ok := module.Flatmap[key]; ok {
		return value
	}

	return nil
}

func (module *ModuleBase) GetKeys() (result []string) {
	module.RWTex.RLock()
	defer module.RWTex.RUnlock()

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
	module.RWTex.RLock()
	defer module.RWTex.RUnlock()

	if module.RefreshC != nil {
		module.RefreshC <- externalSource
	}
}
