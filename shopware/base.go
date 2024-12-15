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

func (s *Shopware) BaseContainer(ctx context.Context) *dagger.Container {
	return dag.Container().From(shopwareBaseImage).WithUser(shopwareUser).
		With(WithBuildDependencies()).
		With(WithBaseEnvironment()).
		With(WithComposerCache(ctx)).
		With(WithNpmCache(ctx)).
		With(WithShopwareSource(s)).
		With(WithComposerDependencies(shopwareProjectRoot)).
		With(WithConfigHMAC(shopwareProjectRoot)).
		With(WithRuntimeVolumes(shopwareProjectRoot))
}

func (s *Shopware) SourceWithVendor(ctx context.Context) *dagger.Directory {
	return s.
		BaseContainer(ctx).
		Directory(shopwareProjectRoot)
}

func WithBuildDependencies() dagger.WithContainerFunc {
	return func(c *dagger.Container) *dagger.Container {
		return c.
			WithFile("/usr/local/bin/composer", dag.Container().From("composer:2").File("/usr/bin/composer"))
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

func WithShopwareSource(s *Shopware, opts ...dagger.ContainerWithMountedDirectoryOpts) dagger.WithContainerFunc {
	if opts == nil {
		opts = []dagger.ContainerWithMountedDirectoryOpts{}
	}

	opts = append(opts, dagger.ContainerWithMountedDirectoryOpts{
		Owner: shopwareUser,
	})

	return func(c *dagger.Container) *dagger.Container {
		return c.
			WithMountedDirectory(shopwareProjectRoot, s.Source, opts...)
	}
}

func WithComposerCache(ctx context.Context, opts ...dagger.ContainerWithMountedCacheOpts) dagger.WithContainerFunc {
	if opts == nil {
		opts = []dagger.ContainerWithMountedCacheOpts{}
	}

	opts = append(opts, dagger.ContainerWithMountedCacheOpts{
		Expand: true,
		Owner:  shopwareUser,
	})

	return func(c *dagger.Container) *dagger.Container {
		if composerHome, err := c.EnvVariable(ctx, "COMPOSER_HOME"); err == nil && composerHome != "" {
			return c.WithMountedCache(composerHome, dag.CacheVolume("composer"), opts...)
		} else {
			return c.WithMountedCache("/root/.composer", dag.CacheVolume("composer"))
		}
	}
}

func WithNpmCache(ctx context.Context, opts ...dagger.ContainerWithMountedCacheOpts) dagger.WithContainerFunc {
	if opts == nil {
		opts = []dagger.ContainerWithMountedCacheOpts{}
	}

	opts = append(opts, dagger.ContainerWithMountedCacheOpts{
		Expand: true,
		Owner:  shopwareUser,
	})

	return func(c *dagger.Container) *dagger.Container {
		if npmCache, err := c.EnvVariable(ctx, "NPM_CONFIG_CACHE"); err == nil && npmCache != "" {
			return c.WithMountedCache(npmCache, dag.CacheVolume("npm"), opts...)
		} else {
			return c.WithMountedCache("/root/.npm", dag.CacheVolume("npm"))
		}
	}
}

func WithRuntimeVolumes(path string, opts ...dagger.ContainerWithMountedCacheOpts) dagger.WithContainerFunc {
	if opts == nil {
		opts = []dagger.ContainerWithMountedCacheOpts{}
	}

	opts = append(opts, dagger.ContainerWithMountedCacheOpts{
		Owner: shopwareUser,
	})

	return func(c *dagger.Container) *dagger.Container {
		return c.
			WithMountedCache(path+"/files", dag.CacheVolume("files"), opts...).
			WithMountedCache(path+"/public/theme", dag.CacheVolume("theme"), opts...).
			WithMountedCache(path+"/public/media", dag.CacheVolume("media"), opts...).
			WithMountedCache(path+"/public/thumbnail", dag.CacheVolume("thumbnail"), opts...).
			WithMountedCache(path+"/public/sitemap", dag.CacheVolume("sitemap"), opts...).
			WithMountedCache(path+"/var/cache", dag.CacheVolume("http_cache"), opts...)
	}
}

func WithConfigHMAC(path string, opts ...dagger.ContainerWithNewFileOpts) dagger.WithContainerFunc {
	config := `
shopware:
  api:
    jwt_key:
      use_app_secret: true
`

	if opts == nil {
		opts = []dagger.ContainerWithNewFileOpts{}
	}

	opts = append(opts, dagger.ContainerWithNewFileOpts{
		Owner: shopwareUser,
	})

	return func(c *dagger.Container) *dagger.Container {
		return c.
			WithNewFile(path+"/config/packages/90-hmac-secret.yaml", config, opts...)
	}
}

func WithComposerDependencies(path string) dagger.WithContainerFunc {
	return func(c *dagger.Container) *dagger.Container {
		return c.
			WithExec([]string{"composer", "install", "-d", path, "-o", "-n", "--no-plugins", "--no-scripts"})
	}
}
