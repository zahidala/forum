package env

import (
	"log"
	"os"
	"strings"
)

// Init initializes the environment variables
func Init() {
	LoadEnv(".env")
}

// LoadEnv loads environment variables from a file
func LoadEnv(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalln("Error reading .env file:", err)
		return err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || line[0] == '#' {
			continue // Skip empty lines and comments
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // Skip lines that do not have exactly one "="
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		os.Setenv(key, value)
	}

	return nil
}
