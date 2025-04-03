package db

type Config struct {
	DebugParser  bool
	DebugStore   bool
	DebugPlanner bool
}

func NewConfig(getEnv func(string) string) *Config {
	return &Config{
		DebugParser:  getEnv("DEBUG_PARSER") == "true",
		DebugStore:   getEnv("DEBUG_STORE") == "true",
		DebugPlanner: getEnv("DEBUG_PLANNER") == "true",
	}
}
