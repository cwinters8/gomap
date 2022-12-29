package utils

import (
	"github.com/joho/godotenv"
)

// load environment variables from path
func Env(path string) error {
	return godotenv.Load(path)
}
