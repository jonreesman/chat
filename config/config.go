package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var config map[string]string
var envs = []string{
	"DB_HOST",
	"DB_PORT",
	"DB_USER",
	"DB_PASSWORD",
	"DB_NAME",
}

func init() {
	config = make(map[string]string)
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Error loading env file: %v", err)
	}
	for _, val := range envs {
		if value := os.Getenv(val); value != "" {
			config[val] = value
		} else {
			panic(fmt.Sprintf("missing value in environment: %s\n please double check your .env file", val))
		}
	}
}

func inSlice(s []string, key string) bool {
	for _, v := range s {
		if v == key {
			return true
		}
	}
	return false
}

func GetConfig(key string) string {
	if val, ok := config[key]; ok {
		return val
	}
	if inSlice(envs, key) {
		value := os.Getenv(key)
		config[key] = value
		return value
	}
	return ""

}
