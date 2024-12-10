package main

import (
	"context"
	"dagger/shopware/internal/dagger"
)

func WithInstall(s *Shopware, ctx context.Context) dagger.WithContainerFunc {
	return func(c *dagger.Container) *dagger.Container {
		return c.
			WithExec([]string{"vendor/bin/shopware-deployment-helper", "run", "--skip-theme-compile", "--skip-assets-install"})
	}
}
