package main

import (
	"context"
	"crypto/rand"
	"dagger/shopware/internal/dagger"
	"encoding/base64"
	"strings"
)

func EnvVariables(envs map[string]string) dagger.WithContainerFunc {
	return func(c *dagger.Container) *dagger.Container {
		for key, val := range envs {
			c = c.WithEnvVariable(key, val)
		}
		return c
	}
}

func WithIntegration(ctx context.Context) dagger.WithContainerFunc {
	bts := make([]byte, 8)
	name := "test"

	if _, err := rand.Read(bts); err == nil {
		name = base64.RawURLEncoding.EncodeToString(bts)
	}

	return func(c *dagger.Container) *dagger.Container {
		env, err := c.
			WithExec([]string{"bin/console", "integration:create", "--admin", name}).
			Stdout(ctx)

		if err != nil {
			return c
		}

		for _, line := range strings.Split(env, "\n") {
			parts := strings.Split(line, "=")
			if len(parts) != 2 {
				continue
			}
			key := parts[0]
			value := parts[1]
			c = c.WithEnvVariable(key, value)
		}

		return c
	}
}
