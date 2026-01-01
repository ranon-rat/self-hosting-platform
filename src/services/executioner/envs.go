package executionerServices

import (
	"os"
	"slices"
	"strings"
)

func LoadLocalEnviroments() []string {
	ourEnvs := os.Environ()
	envs := make([]string, 0, len(ourEnvs))
	for _, env := range ourEnvs {
		if slices.Contains([]string{"PORT", "PASSWORD"}, strings.Split(env, "=")[0]) {
			continue
		}
		envs = append(envs, env)
	}
	return envs
}
