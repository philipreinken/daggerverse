package main

import (
	"context"
	"dagger/shopware/internal/dagger"
)

func (s *Shopware) BasicStack(ctx context.Context) *dagger.Container {
	return s.
		BaseWithDependencies(ctx).
		With(WithCache()).
		With(WithDatabase())
}

func (s *Shopware) Web(ctx context.Context) *dagger.Service {
	return s.
		BasicStack(ctx).
		With(WithFullInstall()).
		WithExposedPort(8000).
		WithExec([]string{"/usr/bin/supervisord", "-c", "/etc/supervisord.conf"}).
		AsService()
}
