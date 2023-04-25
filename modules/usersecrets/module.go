package usersecrets

import (
	"sync"

	"github.com/andreasisnes/go-configuration-manager/modules"
	"github.com/andreasisnes/go-configuration-manager/tools/usersecrets/util"
	"github.com/andreasisnes/goflat"
)

type Options struct {
	modules.Options
}

type Module struct {
	modules.ModuleBase
	options       Options
	waitGroup     sync.WaitGroup
	secrets       map[string]any
	configuration map[string]interface{}
	quitC         chan any
}

func New(options *Options) modules.Module {
	if options == nil {
		options = &Options{}
	}

	return &Module{
		ModuleBase: *modules.NewSourceBase(&options.Options),
		quitC:      make(chan any),
		options:    *options,
		waitGroup:  sync.WaitGroup{},
	}
}

func (module *Module) GetRefreshedValue(key string) any {
	return nil
}

func (module *Module) Load() {
	secrets, err := util.ReadSecrets()
	if err != nil {
		return
	}

	module.secrets = secrets
	module.Flatmap = goflat.Map(module.configuration, &goflat.Options{
		Delimiter: goflat.DefaultDelimiter,
		Fold:      goflat.UpperCaseFold,
	})
}

func (module *Module) Deconstruct() {
}
