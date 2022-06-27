package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

/*
 *Config(): locate secret.env and load database secrets
 */
func Config(key string) string {
	err := godotenv.Load("./dev/secret.env")
	if err != nil {
		log.Fatal("Cant load Secret file")
	}
	return os.Getenv(key)
}
