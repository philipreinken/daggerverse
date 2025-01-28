package main

import (
	"context"
	"dagger/shopware/internal/dagger"
)

const (
	playwrightBaseImage   = "mcr.microsoft.com/playwright:v1.49.1-noble"
	playwrightUser        = "pwuser:pwuser"
	playwrightProjectRoot = shopwareProjectRoot + "/tests/acceptance"
)

func (s *Shopware) PlaywrightContainer(ctx context.Context) *dagger.Container {
	shopware := s.ShopwareService(ctx)

	return dag.Container().From(playwrightBaseImage).
		With(WithShopwareSource(s, dagger.ContainerWithMountedDirectoryOpts{
			Owner: playwrightUser,
		})).
		WithWorkdir(playwrightProjectRoot).
		WithExec([]string{"npm", "ci"}).
		WithExec([]string{"npx", "playwright", "install", "--with-deps", "chromium"}).
		WithServiceBinding("shopware", shopware).
		With(EnvVariables(map[string]string{
			"APP_URL":                 "http://shopware:8000",
			"SHOPWARE_ADMIN_USERNAME": "admin",
			"SHOPWARE_ADMIN_PASSWORD": "shopware",
		})).
		WithUser(playwrightUser)
}
