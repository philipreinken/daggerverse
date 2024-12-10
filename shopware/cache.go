package main

import (
	"context"
	"dagger/shopware/internal/dagger"
)

func WithCache(s *Shopware, ctx context.Context) dagger.WithContainerFunc {
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
