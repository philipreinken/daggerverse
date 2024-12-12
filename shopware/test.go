package main

import (
	"context"
	"dagger/shopware/internal/dagger"
)

func (s *Shopware) Ecs(ctx context.Context) (string, error) {
	return s.
		BaseWithDependencies(ctx).
		With(WithCsFixerCache()).
		WithExec([]string{"composer", "run", "ecs"}).
		Stdout(ctx)
}

func (s *Shopware) LintChangelog(ctx context.Context) (string, error) {
	return s.
		BaseWithDependencies(ctx).
		WithExec([]string{"composer", "run", "lint:changelog"}).
		Stdout(ctx)
}

func (s *Shopware) LintSnippets(ctx context.Context) (string, error) {
	return s.
		BaseWithDependencies(ctx).
		WithExec([]string{"composer", "run", "lint:snippets"}).
		Stdout(ctx)
}

func (s *Shopware) PhpStan(ctx context.Context) (string, error) {
	return s.
		BasicStack(ctx).
		With(WithInstall()).
		WithExec([]string{"composer", "run", "framework:schema:dump"}).
		WithExec([]string{"composer", "run", "phpstan"}).
		Stdout(ctx)
}

func WithCsFixerCache() dagger.WithContainerFunc {
	return func(c *dagger.Container) *dagger.Container {
		return c.
			WithMountedCache(shopwareProjectRoot+"/var/cache/cs_fixer", dag.CacheVolume("cs_fixer"), dagger.ContainerWithMountedCacheOpts{
				Owner: shopwareUser,
			})
	}
}
