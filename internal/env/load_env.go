package env

import (
	"log"
	"os"
	"strings"
)

func LoadDotEnv() {
	file, err := os.Open(".env")
	if err != nil {
		log.Fatalf("couldn't open .env file: %v", err)
	}

	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("couldn't read file info: %v", err)
	}

	bytes := make([]byte, fileInfo.Size())
	_, err = file.Read(bytes)
	if err != nil {
		log.Fatalf("couldn't read bytes of .env file: %v", err)
	}

	envVars := strings.Split(string(bytes), "\n")

	for _, v := range envVars {
		parts := strings.Split(v, "=")
		err := os.Setenv(parts[0], strings.Trim(parts[1], "\""))

		if err != nil {
			log.Fatalf("couldn't set env variable: %v", err)
		}
	}
}