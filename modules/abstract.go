package modules

import "sync"

type Module interface {
	Connect(refreshC chan Module)
	Exists(key string) bool
	Get(key string) any
	GetKeys() []string
	GetOptions() *Options

	GetRefreshedValue(key string) any
	Load()
	Deconstruct()
}

const (
	RefreshAll RefreshPolicy = iota
	RefreshCurrentAndOver
	RefreshCurrentAndUnder
	RefreshCurrent
)

type Options struct {
	Optional        bool
	ReloadOnChange  bool
	Delimiter       string
	SentinelOptions *SentinelOptions
}

type SentinelOptions struct {
	Key           string
	RefreshPolicy RefreshPolicy
}

type ModuleBase struct {
	RefreshC chan Module
	Flatmap  map[string]interface{}
	RWTex    sync.RWMutex
	Options  *Options
}

type RefreshPolicy int
