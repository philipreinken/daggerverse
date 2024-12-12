package main

import (
	"dagger/shopware/internal/dagger"
)

func WithCache() dagger.WithContainerFunc {
	cache := dag.Container().From("valkey/valkey:latest").
		WithExposedPort(6379).
		AsService()

	return func(c *dagger.Container) *dagger.Container {
		return c.
			WithServiceBinding("cache", cache).
			With(EnvVariables(map[string]string{
				"PHP_SESSION_HANDLER":   "redis",
				"PHP_SESSION_SAVE_PATH": "tcp://cache:6379/1",
			}))
	}
}
