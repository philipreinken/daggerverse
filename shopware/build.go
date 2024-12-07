package main

import (
	"context"
	"dagger/shopware/internal/dagger"
)

func (s *Shopware) Build(ctx context.Context) *dagger.Directory {
	return dag.Container().From(shopwareCliImage).
		With(WithBaseEnvironment(s)).
		With(WithConfigHMAC(s, ctx)).
		With(WithShopwareSource(s, ctx)).
		With(WithComposerCache(s, ctx)).
		With(WithDefaultVolumes(s, ctx)).
		WithExec([]string{"shopware-cli", "project", "ci", "--with-dev-dependencies", shopwareProjectRoot}).
		Directory(shopwareProjectRoot)
}

func WithBuildResult(s *Shopware, ctx context.Context, opts ...dagger.ContainerWithMountedDirectoryOpts) dagger.WithContainerFunc {
	return func(c *dagger.Container) *dagger.Container {
		return c.
			WithMountedDirectory(shopwareProjectRoot, s.Build(ctx), opts...).
			WithWorkdir(shopwareProjectRoot)
	}
}
