package config

import (
	"os"
	"strconv"
)

type Config struct {
	Env            string
	Port           string
	MongoURI       string
	JWTSecretKey   string
	EnableMedia    bool
	EnablePosts    bool
	EnableProjects bool
	EnableResume   bool
	EnableUsers    bool
	MediaRoot      string
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
		return fallback
	}
	return fallback
}

func Load() Config {
	cfg := Config{
		Env:            getEnv("APP_ENV", "development"),
		Port:           getEnv("APP_PORT", "8080"),
		MongoURI:       getEnv("MONGO_URI", "mongodb://localhost:27017"),
		JWTSecretKey:   getEnv("JWT_SECRET_KEY", "supersecretkey"),
		EnableMedia:    getEnvBool("ENABLE_MEDIA", false),
		EnablePosts:    getEnvBool("ENABLE_POSTS", false),
		EnableProjects: getEnvBool("ENABLE_PROJECTS", false),
		EnableResume:   getEnvBool("ENABLE_RESUME", false),
		EnableUsers:    getEnvBool("ENABLE_USERS", false),
		MediaRoot:      getEnv("MEDIA_ROOT", "./media"),
	}

	return cfg
}
