package walmanager

type config struct {
	walFilename string
}

func newConfig(opts ...Option) config {
	var c config
	for _, opt := range opts {
		opt(&c)
	}
	return c
}

type Option func(*config)

func WithWALFilename(filename string) Option {
	return func(c *config) {
		c.walFilename = filename
	}
}
