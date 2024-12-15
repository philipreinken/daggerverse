package main

import (
	"context"
	"dagger/shopware/internal/dagger"
)

func (s *Shopware) Ecs(ctx context.Context) *dagger.Container {
	return s.
		DefaultContainer(ctx).
		With(WithCsFixerCache()).
		WithExec([]string{"composer", "run", "ecs"})
}

func (s *Shopware) LintChangelog(ctx context.Context) *dagger.Container {
	return s.
		DefaultContainer(ctx).
		WithExec([]string{"composer", "run", "lint:changelog"})
}

func (s *Shopware) LintSnippets(ctx context.Context) *dagger.Container {
	return s.
		DefaultContainer(ctx).
		WithExec([]string{"composer", "run", "lint:snippets"})
}

func (s *Shopware) Phpstan(ctx context.Context) *dagger.Container {
	return s.
		BasicStack(ctx).
		With(WithInstall()).
		WithExec([]string{"composer", "run", "framework:schema:dump"}).
		WithExec([]string{"composer", "run", "phpstan"})
}

func WithCsFixerCache() dagger.WithContainerFunc {
	return func(c *dagger.Container) *dagger.Container {
		return c.
			WithMountedCache(shopwareProjectRoot+"/var/cache/cs_fixer", dag.CacheVolume("cs_fixer"), dagger.ContainerWithMountedCacheOpts{
				Owner: shopwareUser,
			})
	}
}
