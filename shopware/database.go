package main

import (
	"context"
	"dagger/shopware/internal/dagger"
)

func WithDatabase(s *Shopware, ctx context.Context) dagger.WithContainerFunc {
	database := dag.Container().From("mysql:8.0").
		With(EnvVariables(map[string]string{
			"MYSQL_ALLOW_EMPTY_PASSWORD": "yes",
			"MYSQL_DATABASE":             "shopware",
		})).
		WithMountedTemp("/var/lib/mysql").
		WithExposedPort(3306).
		AsService()

	return func(c *dagger.Container) *dagger.Container {
		return c.
			WithServiceBinding("database", database).
			With(EnvVariables(map[string]string{
				"DATABASE_HOST": "database", // Only used for health checks
				"DATABASE_PORT": "3306",     // Only used for health checks
				"DATABASE_URL":  "mysql://root@database/shopware",
			}))

	}
}
