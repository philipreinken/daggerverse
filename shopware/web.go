package main

import (
	"dagger/shopware/internal/dagger"
)

func WithWebserver(shopwareContainer *dagger.Container) dagger.WithContainerFunc {
	shopware := shopwareContainer.
		WithExposedPort(8000).
		WithExec([]string{"/usr/bin/supervisord", "-c", "/etc/supervisord.conf"}).
		AsService()

	return func(c *dagger.Container) *dagger.Container {
		return c.
			WithServiceBinding("shopware", shopware).
			With(EnvVariables(map[string]string{
				"APP_URL":              "http://shopware:8000",
				"STOREFRONT_PROXY_URL": "http://shopware:8000",
			}))

	}
}
