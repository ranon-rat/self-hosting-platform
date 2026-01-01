package executionerServices

import (
	"os"
	"slices"
)

func LoadLocalEnviroments() []string {
	ourEnvs := os.Environ()
	envs := make([]string, 0, len(ourEnvs))
	for _, env := range ourEnvs {
		if slices.Contains([]string{"PORT", "PASSWORD"}, env) {
			continue
		}
		envs = append(envs, env)
	}
	return envs
}
