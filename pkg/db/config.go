package db

type Config struct {
	DebugScanner bool
	DebugParser  bool
	DebugStore   bool
	DebugPlanner bool
}

func NewConfig(getEnv func(string) string) *Config {
	return &Config{
		DebugScanner: getEnv("DEBUG_SCANNER") == "true",
		DebugParser:  getEnv("DEBUG_PARSER") == "true",
		DebugStore:   getEnv("DEBUG_STORE") == "true",
		DebugPlanner: getEnv("DEBUG_PLANNER") == "true",
	}
}
