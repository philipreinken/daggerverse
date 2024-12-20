package main

import (
	"crypto/rand"
	"dagger/shopware/internal/dagger"
	"encoding/base64"
	"fmt"
)

func WithDatabase() dagger.WithContainerFunc {
	bts := make([]byte, 8)
	name := "shopware"

	if _, err := rand.Read(bts); err == nil {
		name = base64.RawURLEncoding.EncodeToString(bts)
	}

	database := dag.Container().From("mysql:8.0").
		With(EnvVariables(map[string]string{
			"MYSQL_ALLOW_EMPTY_PASSWORD": "yes",
			"MYSQL_DATABASE":             name,
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
				"DATABASE_URL":  fmt.Sprintf("mysql://root@database/%s", name),
			}))

	}
}
