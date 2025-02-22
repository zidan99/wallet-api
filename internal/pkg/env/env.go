package env

import (
	"os"
)

func Get(k string) string {
	return os.Getenv(k)
}

func GetWithDefault(k, d string) string {
	v := os.Getenv(k)
	if v == "" {
		return d
	}

	return v
}
