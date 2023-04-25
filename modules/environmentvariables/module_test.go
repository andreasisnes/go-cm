package environmentvariables

import (
	"os"
	"testing"
	"time"

	configurationmanager "github.com/andreasisnes/go-configuration-manager"
	"github.com/andreasisnes/go-configuration-manager/modules"
	"github.com/stretchr/testify/assert"
)

func newConfig(option *Options) configurationmanager.Configuration {
	return configurationmanager.New().
		Add(New(option)).
		Build()
}

func TestGet(t *testing.T) {
	key := "TestGet"
	expected := "TEST_VALUE"
	os.Setenv(key, expected)
	config := newConfig(nil)

	assert.Equal(t, expected, config.Get(key, nil))
}

func TestGetWithReloadOnChange(t *testing.T) {
	config := newConfig(&Options{
		RefreshInterval: time.Second,
		Options: modules.Options{
			ReloadOnChange: true,
		},
	})

	key := "TestGetWithReloadOnChange"
	expected := "TEST_VALUE"
	os.Setenv(key, expected)
	time.Sleep(time.Second * 2)
	result := config.Get(key, nil)

	assert.Equal(t, expected, result)
	config.Deconstruct()
}

func TestGetWithSentinel(t *testing.T) {
	t.Parallel()
	config := newConfig(&Options{
		RefreshInterval: time.Second,
		Options: modules.Options{
			ReloadOnChange: false,
			SentinelOptions: &modules.SentinelOptions{
				Key:           "TestGetWithSentinel",
				RefreshPolicy: modules.RefreshCurrent,
			},
		},
	})

	key := "TestGetWithSentinel"
	expected := "TEST_VALUE"
	os.Setenv(key, expected)
	time.Sleep(time.Second * 2)
	result := config.Get(key, nil)

	assert.Equal(t, expected, result)
	config.Deconstruct()
}

func TestGetNilWithSentinel(t *testing.T) {
	t.Parallel()
	config := newConfig(&Options{
		RefreshInterval: time.Second,
		Options: modules.Options{
			ReloadOnChange: false,
			SentinelOptions: &modules.SentinelOptions{
				Key:           "TestGetNilWithSentinelUnkownKey",
				RefreshPolicy: modules.RefreshCurrent,
			},
		},
	})

	key := "TestGetNilWithSentinel"
	notExpected := "TEST_VALUE"
	os.Setenv(key, notExpected)
	time.Sleep(time.Second * 2)
	result := config.Get(key, nil)

	assert.Nil(t, result)
	config.Deconstruct()
}
