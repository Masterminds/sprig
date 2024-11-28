package sprig

import (
	"fmt"
	"os"
)

func mustEnv(env string) (string, error){
	v, ok := os.LookupEnv(env)
	if !ok {
		return v, fmt.Errorf("env var \"%s\" does not exist", env)
	}

	return v, nil
}