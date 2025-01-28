package main

import (
	"dagger/shopware/internal/dagger"
	"fmt"
	"strings"
)

// FIXME: This will only work in combination with a minio/S3 service, since the `public` directory is not synced between containers
func WithWebserver() dagger.WithContainerFunc {
	return func(c *dagger.Container) *dagger.Container {
		shopware := c.
			WithExposedPort(8000).
			WithExec([]string{"/usr/bin/supervisord", "-c", "/etc/supervisord.conf", "--nodaemon"}).
			AsService()

		return c.
			WithServiceBinding("shopware", shopware).
			With(EnvVariables(map[string]string{
				"APP_URL":              "http://shopware:8000",
				"STOREFRONT_PROXY_URL": "http://shopware:8000",
			}))

	}
}

func WithWebServerAndExec(exec []string) dagger.WithContainerFunc {
	cmdSrv := []string{"/usr/bin/supervisord", "-c", "/etc/supervisord.conf"}

	return func(c *dagger.Container) *dagger.Container {
		return c.
			WithExec([]string{
				"bash", "-c", fmt.Sprintf("%s && %s", strings.Join(cmdSrv, " "), strings.Join(exec, " ")),
			})
	}
}
