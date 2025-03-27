package db

type Config struct {
	DebugParser bool
	DebugStore  bool
}

func NewConfig(getEnv func(string) string) *Config {
	return &Config{
		DebugParser: getEnv("DEBUG_PARSER") == "true",
		DebugStore:  getEnv("DEBUG_STORE") == "true",
	}
}
