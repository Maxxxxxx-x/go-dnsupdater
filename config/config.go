package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type AuthConfig struct {
	Token string
}

type Config struct {
	Auth    AuthConfig
	Domains []string
}

func missingEnv(envName string) string {
	return fmt.Sprintf("Missing ENV variable: %s", envName)
}

func getAuthConfig() AuthConfig {
	token, found := os.LookupEnv("API_TOKEN")
	if !found {
		log.Fatalln(missingEnv("API_TOKEN"))
	}

	return AuthConfig{
		Token: token,
	}
}

func getDomains() []string {
	domains, found := os.LookupEnv("DOMAINS")
	if !found {
		log.Fatalln(missingEnv("DOMAINS"))
	}
	return strings.Split(domains, ",")
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to parse .env file: %s\n", err.Error())
	}

	return Config{
		Auth:    getAuthConfig(),
		Domains: getDomains(),
	}
}
