package main

import (
	"dagger/shopware/internal/dagger"
)

const (
	// The base image to use for the shopware container
	shopwareBaseImage = "ghcr.io/shopware/docker-base:8.3-nginx-otel"
	shopwareCliImage  = "ghcr.io/friendsofshopware/shopware-cli:latest-php-8.3"
	// Non-Root user to use for the shopware container
	shopwareUser = "www-data:www-data"
	// Workdir for the shopware container (PROJECT_ROOT)
	shopwareProjectRoot = "/var/www/html"
)

func (s *Shopware) Base() *dagger.Container {
	return dag.Container().From(shopwareBaseImage).
		With(WithBaseEnvironment(s)).
		With(WithBuildResult(s, dagger.ContainerWithMountedDirectoryOpts{
			Owner: shopwareUser,
		})).
		With(WithDefaultVolumes(s, dagger.ContainerWithMountedCacheOpts{
			Owner: shopwareUser,
		}))
}

func (s *Shopware) Build() *dagger.Directory {
	return dag.Container().From(shopwareCliImage).
		With(WithBaseEnvironment(s)).
		With(WithShopwareSource(s)).
		With(WithDefaultVolumes(s)).
		WithExec([]string{"shopware-cli", "project", "ci", "--with-dev-dependencies", shopwareProjectRoot}).
		Directory(shopwareProjectRoot)
}

func WithBaseEnvironment(s *Shopware) dagger.WithContainerFunc {
	return func(c *dagger.Container) *dagger.Container {
		return c.
			With(EnvVariables(map[string]string{
				"APP_ENV":                              "prod",
				"APP_DEBUG":                            "1",
				"APP_URL":                              "http://127.0.0.1:8000",
				"BLUE_GREEN_DEPLOYMENT":                "0",
				"COMPOSER_PLUGIN_LOADER":               "1",
				"COMPOSER_ROOT_VERSION":                "6.6.9999999-dev",
				"DATABASE_URL":                         "mysql://null",
				"MAILER_DSN":                           "null://localhost",
				"SHOPWARE_ES_ENABLED":                  "0",
				"SHOPWARE_ES_INDEXING_ENABLED":         "0",
				"SHOPWARE_HTTP_CACHE_ENABLED":          "0",
				"MESSENGER_TRANSPORT_DSN":              "null://localhost",
				"MESSENGER_TRANSPORT_LOW_PRIORITY_DSN": "null://localhost",
				"MESSENGER_TRANSPORT_FAILURE_DSN":      "null://localhost",
			}))
	}
}

func WithShopwareSource(s *Shopware, opts ...dagger.ContainerWithMountedDirectoryOpts) dagger.WithContainerFunc {
	return func(c *dagger.Container) *dagger.Container {
		return c.
			WithMountedDirectory(shopwareProjectRoot, s.Source, opts...).
			WithWorkdir(shopwareProjectRoot)
	}
}

func WithDefaultVolumes(s *Shopware, opts ...dagger.ContainerWithMountedCacheOpts) dagger.WithContainerFunc {
	return func(c *dagger.Container) *dagger.Container {
		return c.
			WithMountedCache(shopwareProjectRoot+"/files", dag.CacheVolume("files"), opts...).
			WithMountedCache(shopwareProjectRoot+"/public/theme", dag.CacheVolume("theme"), opts...).
			WithMountedCache(shopwareProjectRoot+"/public/media", dag.CacheVolume("media"), opts...).
			WithMountedCache(shopwareProjectRoot+"/public/thumbnail", dag.CacheVolume("thumbnail"), opts...).
			WithMountedCache(shopwareProjectRoot+"/public/sitemap", dag.CacheVolume("sitemap"), opts...).
			WithMountedCache(shopwareProjectRoot+"/var/cache", dag.CacheVolume("http_cache"), opts...)
	}
}

func WithBuildResult(s *Shopware, opts ...dagger.ContainerWithMountedDirectoryOpts) dagger.WithContainerFunc {
	return func(c *dagger.Container) *dagger.Container {
		return c.
			WithMountedDirectory(shopwareProjectRoot, s.Build(), opts...).
			WithWorkdir(shopwareProjectRoot)
	}
}

func EnvVariables(envs map[string]string) dagger.WithContainerFunc {
	return func(c *dagger.Container) *dagger.Container {
		for key, val := range envs {
			c = c.WithEnvVariable(key, val)
		}
		return c
	}
}
