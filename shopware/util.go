package main

import (
	"dagger/shopware/internal/dagger"
)

func EnvVariables(envs map[string]string) dagger.WithContainerFunc {
	return func(c *dagger.Container) *dagger.Container {
		for key, val := range envs {
			c = c.WithEnvVariable(key, val)
		}
		return c
	}
}
