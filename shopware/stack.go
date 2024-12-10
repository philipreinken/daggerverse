package main

import (
	"context"
	"dagger/shopware/internal/dagger"
)

func (s *Shopware) BasicStack(ctx context.Context) *dagger.Container {
	return s.
		Base(ctx).
		With(WithCache(s, ctx)).
		With(WithDatabase(s, ctx)).
		With(WithInstall(s, ctx))
}

func (s *Shopware) Web(ctx context.Context) *dagger.Service {
	return s.
		BasicStack(ctx).
		WithExposedPort(8000).
		WithExec([]string{"/usr/bin/supervisord", "-c", "/etc/supervisord.conf"}).
		AsService()
}
