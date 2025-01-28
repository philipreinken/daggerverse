package main

import (
	"context"
)

func (s *Shopware) Playwright(ctx context.Context) (string, error) {
	return s.
		PlaywrightContainer(ctx).
		WithExposedPort(9323).
		WithExec([]string{"npx", "playwright", "test", "--project", "Platform"}).
		Stdout(ctx)
}
