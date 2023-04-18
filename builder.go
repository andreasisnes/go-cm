package configurationmanager

import (
	"github.com/andreasisnes/go-configuration-manager/modules"
)

type Builder interface {
	Clear()
	Add(module modules.Module) Builder
	Sources() []modules.Module
}

type builder[T any] struct {
	sources []modules.Module
}

func New[T any]() Builder {
	return &builder[T]{
		sources: make([]modules.Module, 0),
	}
}

func (this *builder[T]) Clear() {
	this.sources = make([]modules.Module, 0)
}

func (this *builder[T]) Add(source modules.Module) Builder {
	this.sources = append(this.sources, source)
	return this
}

func (this *builder[T]) Sources() []modules.Module {
	return this.sources
}

func (this *builder[T]) Build() Configuration {
	for _, c := range this.sources {
		c.Load()
	}

	return newConfiguration(this.Sources())
}
