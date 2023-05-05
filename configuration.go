package configurationmanager

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/andreasisnes/go-configuration-manager/modules"
)

var (
	ErrKeyNotFound      = errors.New("key not found")
	ErrPointerNotPassed = errors.New("pointer not passed")
)

type Configuration interface {
	Get(key string, value any) any
	Deconstruct() Configuration
	Refresh() (isRefreshed bool)
	Unmarshal(value any) error
}

type configuration struct {
	refreshC chan modules.Module
	quitC    chan struct{}

	waitgroup sync.WaitGroup
	modules   []modules.Module
	delimiter string
}

func newConfiguration(sources []modules.Module) Configuration {
	config := &configuration{
		waitgroup: sync.WaitGroup{},
		modules:   sources,
		delimiter: ".",
		refreshC:  make(chan modules.Module),
		quitC:     make(chan struct{}),
	}

	for _, source := range config.modules {
		source.Load()
		source.Connect(config.refreshC)
	}

	config.Refresh()
	go config.autoRefresh()

	return config
}

// Get
func (c *configuration) Get(key string, out any) any {
	key = strings.ToUpper(key)
	for idx := range c.modules {
		source := c.modules[len(c.modules)-1-idx]
		if source.Exists(key) {
			value := source.Get(key)
			if out == nil {
				return value
			}

			CastValue(value, &out)
		}
	}

	return out
}

// Refresh
func (c *configuration) Refresh() (successfullyRefreshed bool) {
	defer func() {
		if r := recover(); r != nil {
			successfullyRefreshed = false
			fmt.Println("Recovered. Error:\n", r)
		} else {
			successfullyRefreshed = true
		}
	}()

	wg := sync.WaitGroup{}
	for _, source := range c.modules {
		wg.Add(1)
		go func(sourceArg modules.Module) {
			defer wg.Done()
			sourceArg.Load()
		}(source)
	}
	wg.Wait()

	return successfullyRefreshed
}

// Unmarshal
func (c *configuration) Unmarshal(value any) error {
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return ErrPointerNotPassed
	}

	keys := make(map[string]modules.Module)
	for _, source := range c.modules {
		for _, key := range source.GetKeys() {
			keys[key] = source
		}
	}

	flat := make(map[string]interface{})
	for key, source := range keys {
		flat[key] = source.Get(key)
	}

	return nil
}

// Deconstruct
func (c *configuration) Deconstruct() Configuration {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from error:\n", r)
		}
	}()

	c.quitC <- struct{}{}
	wg := sync.WaitGroup{}
	for _, source := range c.modules {
		wg.Add(1)
		go func(sourceArg modules.Module) {
			defer wg.Done()
			sourceArg.Deconstruct()
		}(source)
	}
	wg.Wait()
	c.waitgroup.Wait()

	return c
}

func (c *configuration) autoRefresh() {
	c.waitgroup.Add(1)
	defer c.waitgroup.Done()

	for {
		select {
		case source := <-c.refreshC:
			if source.GetOptions().ReloadOnChange {
				source.Load()
			}
			if source.GetOptions().SentinelOptions != nil {
				c.loadSentinel(source)
			}
		case <-c.quitC:
			return
		}
	}
}

func (c *configuration) loadSentinel(source modules.Module) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from error:\n", r)
		}
	}()

	key := source.GetOptions().SentinelOptions.Key
	if reflect.DeepEqual(source.Get(key), source.GetRefreshedValue(key)) {
		return
	}

	switch source.GetOptions().SentinelOptions.RefreshPolicy {
	case modules.RefreshAll:
		c.Refresh()
	case modules.RefreshCurrent:
		source.Load()
	case modules.RefreshCurrentAndOver:
		c.refreshCurrentAndAbove(source)
	case modules.RefreshCurrentAndUnder:
		c.refreshCurrentAndUnder(source)
	}
}

func (c *configuration) refreshCurrentAndAbove(source modules.Module) {
	wg := sync.WaitGroup{}
	isAbove := false
	for _, module := range c.modules {
		if module == source {
			isAbove = true
		}
		if isAbove {
			wg.Add(1)
			go func(sourceArg modules.Module) {
				defer wg.Done()
				sourceArg.Load()
			}(source)
		}
	}
	wg.Wait()
}

func (c *configuration) refreshCurrentAndUnder(source modules.Module) {
	wg := sync.WaitGroup{}
	isUnder := true
	for _, s := range c.modules {
		if isUnder {
			wg.Add(1)
			go func(sourceArg modules.Module) {
				defer wg.Done()
				sourceArg.Load()
			}(source)
		}
		if s == source {
			isUnder = false
		}
	}
	wg.Wait()
}
