package main

import (
	"context"
	"dagger/shopware/internal/dagger"
	"encoding/base64"
)

const (
	// The base image to use for the shopware container
	shopwareBaseImage = "ghcr.io/shopware/docker-base:8.3-nginx"
	// Non-Root user to use for the shopware container
	shopwareUser = "www-data:www-data"
	// Workdir for the shopware container (PROJECT_ROOT)
	shopwareProjectRoot = "/var/www/html"
)

func (s *Shopware) Base(ctx context.Context) *dagger.Container {
	return dag.Container().From(shopwareBaseImage).
		With(WithBuildDependencies()).
		With(WithBaseEnvironment()).
		With(WithShopwareSource(s, ctx)).
		With(WithDefaultVolumes(s, ctx)).
		With(WithComposerCache(s, ctx)).
		With(WithNpmCache(s, ctx)).
		With(WithConfigHMAC(s, ctx)).
		WithWorkdir(shopwareProjectRoot)
}

func (s *Shopware) BaseWithDependencies(ctx context.Context) *dagger.Container {
	return s.
		Base(ctx).
		With(WithDependencies())
}

func WithBuildDependencies() dagger.WithContainerFunc {
	return func(c *dagger.Container) *dagger.Container {
		return c.
			WithUser("root:root").
			WithFile("/usr/local/bin/composer", c.From("composer:2").File("/usr/bin/composer")).
			WithExec([]string{"apk", "add", "--no-cache", "nodejs", "npm", "bash"}).
			WithUser(shopwareUser)
	}
}

func WithBaseEnvironment() dagger.WithContainerFunc {
	return func(c *dagger.Container) *dagger.Container {
		return c.
			With(EnvVariables(map[string]string{
				"APP_ENV":                      "dev",
				"APP_DEBUG":                    "1",
				"APP_URL":                      "http://127.0.0.1:8000",
				"APP_SECRET":                   base64.RawURLEncoding.EncodeToString([]byte("0000AAAAshopware0000AAAAshopware")),
				"BLUE_GREEN_DEPLOYMENT":        "0",
				"COMPOSER_PLUGIN_LOADER":       "1",
				"COMPOSER_ROOT_VERSION":        "6.6.9999999-dev",
				"DATABASE_URL":                 "mysql://null",
				"MAILER_DSN":                   "null://null",
				"SHOPWARE_ES_ENABLED":          "0",
				"SHOPWARE_ES_INDEXING_ENABLED": "0",
				"SHOPWARE_HTTP_CACHE_ENABLED":  "0",
				"COMPOSER_HOME":                "/tmp/composer",
				"NPM_CONFIG_CACHE":             "/tmp/npm",
			}))
	}
}

func WithShopwareSource(s *Shopware, ctx context.Context, opts ...dagger.ContainerWithMountedDirectoryOpts) dagger.WithContainerFunc {
	opts = append(opts, dagger.ContainerWithMountedDirectoryOpts{
		Owner: shopwareUser,
	})

	return func(c *dagger.Container) *dagger.Container {
		return c.
			WithMountedDirectory(shopwareProjectRoot, s.Source, opts...)
	}
}

func WithComposerCache(s *Shopware, ctx context.Context, opts ...dagger.ContainerWithMountedCacheOpts) dagger.WithContainerFunc {
	return func(c *dagger.Container) *dagger.Container {
		opts = append(opts, dagger.ContainerWithMountedCacheOpts{
			Expand: true,
			Owner:  shopwareUser,
		})

		if composerHome, err := c.EnvVariable(ctx, "COMPOSER_HOME"); err == nil && composerHome != "" {
			return c.WithMountedCache(composerHome, dag.CacheVolume("composer"), opts...)
		} else {
			return c.WithMountedCache("/root/.composer", dag.CacheVolume("composer"))
		}
	}
}

func WithNpmCache(s *Shopware, ctx context.Context, opts ...dagger.ContainerWithMountedCacheOpts) dagger.WithContainerFunc {
	return func(c *dagger.Container) *dagger.Container {
		opts = append(opts, dagger.ContainerWithMountedCacheOpts{
			Expand: true,
			Owner:  shopwareUser,
		})

		if npmCache, err := c.EnvVariable(ctx, "NPM_CONFIG_CACHE"); err == nil && npmCache != "" {
			return c.WithMountedCache(npmCache, dag.CacheVolume("npm"), opts...)
		} else {
			return c.WithMountedCache("/root/.npm", dag.CacheVolume("npm"))
		}
	}
}

func WithDefaultVolumes(s *Shopware, ctx context.Context, opts ...dagger.ContainerWithMountedCacheOpts) dagger.WithContainerFunc {
	opts = append(opts, dagger.ContainerWithMountedCacheOpts{
		Owner: shopwareUser,
	})

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

func WithConfigHMAC(s *Shopware, ctx context.Context, opts ...dagger.ContainerWithNewFileOpts) dagger.WithContainerFunc {
	config := `
shopware:
  api:
    jwt_key:
      use_app_secret: true
`

	opts = append(opts, dagger.ContainerWithNewFileOpts{
		Owner: shopwareUser,
	})

	return func(c *dagger.Container) *dagger.Container {
		return c.
			WithNewFile(shopwareProjectRoot+"/config/packages/90-hmac-secret.yaml", config, opts...)
	}
}

func WithDependencies() dagger.WithContainerFunc {
	return func(c *dagger.Container) *dagger.Container {
		return c.
			WithExec([]string{"composer", "install", "-o", "-n"})
	}
}
