package env

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

func loadOptionalEnvFile() {
	// It's okay if the .env file doesn't exist, we can proceed with existing env vars
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		panic(err)
	}
}
