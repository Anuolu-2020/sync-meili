package env

import (
	"log"
	"os"
)

func GetEnv(name string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		log.Fatalf("%s env variable does not exist", name)
	}

	return value
}
