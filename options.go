package kotlingo

type Config struct {
	logger Logger
}

func defaultCommonConfig() Config {
	return Config{
		logger: EmptyLogger{},
	}
}

type Option func(cfg *Config)

func WithLogger(logger Logger) Option {
	return func(cfg *Config) {
		cfg.logger = logger
	}
}
