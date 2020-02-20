package main

import (
	"encoding/json"
	"os"

	"gopkg.in/yaml.v2"
)

// JSONizeStruct func
func JSONizeStruct(t interface{}) string {
	b, err := json.Marshal(&t)
	if err != nil {
		return ""
	}
	return string(b)
}

// YAMLizeStruct func
func YAMLizeStruct(t interface{}) string {
	b, err := yaml.Marshal(&t)
	if err != nil {
		return ""
	}
	return string(b)
}

// Get env with fallback
func getEnv(key string, fallback string) string {
	env := os.Getenv(key)
	if len(env) == 0 {
		env = fallback
	}
	return env
}
