package main

import (
	"context"
	"dagger/shopware/internal/dagger"
	"encoding/base64"
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

func (s *Shopware) Base(ctx context.Context) *dagger.Container {
	return dag.Container().From(shopwareBaseImage).
		With(WithBaseEnvironment(s)).
		With(WithBuildResult(s, ctx, dagger.ContainerWithMountedDirectoryOpts{
			Owner: shopwareUser,
		})).
		With(WithDefaultVolumes(s, ctx, dagger.ContainerWithMountedCacheOpts{
			Owner: shopwareUser,
		}))
}

func WithBaseEnvironment(s *Shopware) dagger.WithContainerFunc {
	return func(c *dagger.Container) *dagger.Container {
		return c.
			With(EnvVariables(map[string]string{
				"APP_ENV":                              "prod",
				"APP_DEBUG":                            "1",
				"APP_URL":                              "http://127.0.0.1:8000",
				"APP_SECRET":                           base64.RawURLEncoding.EncodeToString([]byte("0000AAAAshopware0000AAAAshopware")),
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

func WithShopwareSource(s *Shopware, ctx context.Context, opts ...dagger.ContainerWithMountedDirectoryOpts) dagger.WithContainerFunc {
	return func(c *dagger.Container) *dagger.Container {
		return c.
			WithMountedDirectory(shopwareProjectRoot, s.Source, opts...).
			WithWorkdir(shopwareProjectRoot)
	}
}

func WithComposerCache(s *Shopware, ctx context.Context) dagger.WithContainerFunc {
	return func(c *dagger.Container) *dagger.Container {
		if composerHome, err := c.EnvVariable(ctx, "COMPOSER_HOME"); err == nil && composerHome != "" {
			return c.WithMountedCache(composerHome, dag.CacheVolume("composer"), dagger.ContainerWithMountedCacheOpts{
				Expand: true,
			})
		} else {
			return c.WithMountedCache("/root/.composer", dag.CacheVolume("composer"))
		}
	}
}

func WithDefaultVolumes(s *Shopware, ctx context.Context, opts ...dagger.ContainerWithMountedCacheOpts) dagger.WithContainerFunc {
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

func WithConfigHMAC(s *Shopware, ctx context.Context) dagger.WithContainerFunc {
	config := `
shopware:
  api:
    jwt_key:
      use_app_secret: true
`
	return func(c *dagger.Container) *dagger.Container {
		return c.
			WithNewFile(shopwareProjectRoot+"/config/90-hmac-secret.yaml", config)
	}
}
