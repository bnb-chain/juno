package config

type Config struct {
	ServiceName string `yaml:"service_name"`
	RootDir     string `yaml:"root_dir"`
	Level       string `yaml:"level"`
}

// DefaultLogConfig returns the default Config instance
func DefaultLogConfig() Config {
	return Config{
		RootDir: "./logs",
		Level:   "debug",
	}
}
