package configurationmanager

func NewDefaultOptions(optsfn ...func(options *Options)) *Options {
	options := &Options{
		Delimiter: ":",
	}

	for _, fn := range optsfn {
		fn(options)
	}

	return options
}

type Options struct {
	Delimiter string
}
