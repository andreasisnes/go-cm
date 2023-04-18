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

type Configuration[T any] interface {
	Bind(T)
	Create() T
	Get(key string, value interface{}) interface{}
	Deconstruct() Configuration[T]
	Refresh() (successfullyRefreshed bool)
	Unmarshal(value interface{}) error
}

type configuration[T any] struct {
	RefreshC chan modules.Module
	QuitC    chan struct{}

	Waitgroup sync.WaitGroup
	Sources   []modules.Module
	Delimiter string
}

func newConfiguration[T any](sources []modules.Module) Configuration[T] {
	config := &configuration[T]{
		Waitgroup: sync.WaitGroup{},
		Sources:   sources,
		Delimiter: ".",
		RefreshC:  make(chan modules.Module),
		QuitC:     make(chan struct{}),
	}

	for _, source := range config.Sources {
		source.Connect(config.RefreshC)
	}

	config.Refresh()
	go config.autoRefresh()

	return config
}

func (c *configuration) Get(key string, result interface{}) interface{} {
	key = strings.ToUpper(key)
	for idx := range c.Sources {
		source := c.Sources[len(c.Sources)-1-idx]
		if source.Exists(key) {
			value := source.Get(key)
			if result == nil {
				return value
			}

			return CastAndTryAssignValue(value, result)
		}
	}

	return result
}

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
	for _, source := range c.Sources {
		wg.Add(1)
		go func(sourceArg modules.Module) {
			defer wg.Done()
			sourceArg.Load()
		}(source)
	}
	wg.Wait()

	return successfullyRefreshed
}

func (c *configuration) Unmarshal(value interface{}) error {
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return ErrPointerNotPassed
	}

	keys := make(map[string]modules.Module)
	for _, source := range c.Sources {
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

func (c *configuration) Deconstruct() Configuration {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from error:\n", r)
		}
	}()

	c.QuitC <- struct{}{}
	wg := sync.WaitGroup{}
	for _, source := range c.Sources {
		wg.Add(1)
		go func(sourceArg modules.Module) {
			defer wg.Done()
			sourceArg.Deconstruct()
		}(source)
	}
	wg.Wait()
	c.Waitgroup.Wait()

	return c
}

func (c *configuration) autoRefresh() {
	c.Waitgroup.Add(1)
	defer c.Waitgroup.Done()

	for {
		select {
		case source := <-c.RefreshC:
			if source.GetOptions().ReloadOnChange {
				source.Load()
			}
			if source.GetOptions().SentinelOptions != nil {
				c.loadSentinel(source)
			}
		case <-c.QuitC:
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
	for _, s := range c.Sources {
		if s == source {
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
	for _, s := range c.Sources {
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
