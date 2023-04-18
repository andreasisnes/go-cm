package gonfigenvironmentvariables

import (
	"os"
	"strings"
	"sync"
	"time"

	"github.com/andreasisnes/go-configuration-manager/modules"
)

const (
	DefaultDelimiter       = "__"
	DefaultRefreshInterval = time.Second * 10
)

type Options struct {
	modules.Options
	Prefix          string
	Delimiter       string
	RefreshInterval time.Duration
}

type Module struct {
	modules.ModuleBase
	EnvOptions Options
	WaitGroup  sync.WaitGroup
	QuitC      chan interface{}
}

func New(options *Options) modules.Module {
	if options == nil {
		options = &Options{
			Delimiter:       DefaultDelimiter,
			RefreshInterval: DefaultRefreshInterval,
		}
	}

	if options.Delimiter == "" {
		options.Delimiter = DefaultDelimiter
	}

	module := &Module{
		ModuleBase: *modules.NewSourceBase(&options.Options),
		QuitC:      make(chan interface{}),
		EnvOptions: *options,
		WaitGroup:  sync.WaitGroup{},
	}

	if module.Options.SentinelOptions != nil || module.Options.ReloadOnChange {
		go module.watcher()
	}

	return module
}

func (e *Module) Load() {
	e.RWTex.Lock()
	defer e.RWTex.Unlock()
	for _, variable := range os.Environ() {
		keyIdx := strings.Index(variable, "=")
		e.Flatmap[variable[:keyIdx]] = variable[keyIdx+1:]
		e.Flatmap[e.formatKey(variable[:keyIdx])] = variable[keyIdx+1:]
	}
}

func (e *Module) GetRefreshedValue(key string) interface{} {
	key = e.formatKey(key)
	for _, variable := range os.Environ() {
		keyIdx := strings.Index(variable, "=")
		if e.formatKey(variable[:keyIdx]) == key {
			return variable[keyIdx+1:]
		}
	}

	return nil
}

func (e *Module) Deconstruct() {
	e.QuitC <- struct{}{}
}

func (e *Module) watcher() {
	e.WaitGroup.Add(1)
	defer e.WaitGroup.Done()
	timer := time.NewTimer(e.EnvOptions.RefreshInterval)
	for {
		select {
		case <-e.QuitC:
			return
		case <-timer.C:
			e.RWTex.RLock()
			for _, variable := range os.Environ() {
				keyIdx := strings.Index(variable, "=")
				key := variable[:keyIdx]
				if val, ok := e.Flatmap[key]; !ok || val != variable[keyIdx+1:] {
					e.NotifyDirtyness(e)
					break
				}

				formattedKey := e.formatKey(key)
				if val, ok := e.Flatmap[formattedKey]; !ok || val != variable[keyIdx+1:] {
					e.NotifyDirtyness(e)
					break
				}
			}
			e.RWTex.RUnlock()
		}
		timer.Reset(e.EnvOptions.RefreshInterval)
	}
}

func (e *Module) formatKey(key string) string {
	key = strings.ToUpper(strings.ReplaceAll(key, e.EnvOptions.Delimiter, "."))
	prefix := strings.ToUpper(e.EnvOptions.Prefix)
	return strings.TrimPrefix(key, prefix)
}
