package utils

import (
	"strconv"
	"strings"
)

// UpdateENVWorkerCounts ...
func UpdateENVWorkerCounts(envs []string, count int) []string {
	newEnvs := make([]string, len(envs))
	for i, env := range envs {
		if strings.Contains(env, "WorkerCount") {
			newEnv := "WorkerCount=" + strconv.Itoa(count)
			newEnvs[i] = newEnv
		} else {
			newEnvs[i] = envs[i]
		}
	}
	return newEnvs
}

// ContainsString ...
func ContainsString(list []string, target string) bool {
	for _, value := range list {
		if value == target {
			return true
		}
	}
	return false
}
