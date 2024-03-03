package sys

import "os"

func EnvStr(k, def string) string {
	str := os.Getenv(k)
	if len(str) == 0 {
		return def
	}

	return str
}
