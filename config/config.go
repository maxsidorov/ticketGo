package config

type Config struct {
	DatabasePath string `mapstructure:"DB_PATH"`
	Port         string `mapstructure:"PORT"`
}

func Load() *Config {
	return &Config{
		DatabasePath: getEnv("DB_PATH", "./afisha.db"),
		Port:         getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
