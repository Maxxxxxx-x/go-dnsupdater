package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)


type AuthConfig struct {
    Key string
    Email string
    Token string
}

type DnsConfig struct {
    ZoneId string
    DNSRecordId string
}

type Config struct {
    Auth AuthConfig
    DNS DnsConfig
}


func missingEnv(envName string) string {
    return fmt.Sprintf("Missing ENV variable: %s", envName)
}


func getAuthConfig() AuthConfig {
    key, found := os.LookupEnv("API_KEY")
    if !found {
        log.Fatalln(missingEnv("API_KEY"))
    }

    email, found := os.LookupEnv("API_EMAIL")
    if !found {
        log.Fatalln(missingEnv("API_EMAIL"))
    }

    token, found := os.LookupEnv("API_TOKEN")
    if !found {
        log.Fatalln(missingEnv("API_TOKEN"))
    }

    return AuthConfig{
        Key: key,
        Email: email,
        Token: token,
    }
}

func getDnsConfig() DnsConfig {
    zoneId, found := os.LookupEnv("ZONE_ID")
    if !found {
        log.Fatalln(missingEnv("ZONE_ID"))
    }

    dnsRecordId, found := os.LookupEnv("DNS_RECORD_ID")
    if !found {
        log.Fatalln(missingEnv("DNS_RECORD_ID"))
    }

    return DnsConfig{
        ZoneId: zoneId,
        DNSRecordId: dnsRecordId,
    }
}


func LoadConfig() Config {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Failed to parse .env file: %s\n", err.Error())
    }


    authConfig := getAuthConfig()
    dnsConfig := getDnsConfig()

    return Config{
        Auth: authConfig,
        DNS: dnsConfig,
    }
}
