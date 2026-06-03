package utils

import (
	"os"

	"github.com/joho/godotenv"
)

// load environment variables from path
func Env(path string) error {
	if err := godotenv.Load(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
