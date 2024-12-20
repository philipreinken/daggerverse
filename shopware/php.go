package main

import (
	"context"
	"dagger/shopware/internal/dagger"
	"fmt"
	"golang.org/x/sync/errgroup"
	"strings"
)

var (
	integrationTestPaths = []string{
		"tests/integration/Administration",
		"tests/integration/Core/Checkout",
		"tests/integration/Core/Content",
		"tests/integration/Core/Framework",
		"tests/integration/Core/Installer",
		"tests/integration/Core/Maintenance",
		"tests/integration/Core/System",
		"tests/integration/Elasticsearch",
		"tests/integration/Storefront",
	}
)

func (s *Shopware) Ecs(ctx context.Context) (string, error) {
	return s.
		BaseContainer(ctx).
		With(WithTestEnvironment()).
		With(WithCsFixerCache()).
		WithExec([]string{"composer", "run", "ecs"}).
		Stdout(ctx)
}

func (s *Shopware) LintChangelog(ctx context.Context) (string, error) {
	return s.
		BaseContainer(ctx).
		With(WithTestEnvironment()).
		WithExec([]string{"composer", "run", "lint:changelog"}).
		Stdout(ctx)
}

func (s *Shopware) LintSnippets(ctx context.Context) (string, error) {
	return s.
		BaseContainer(ctx).
		With(WithTestEnvironment()).
		WithExec([]string{"composer", "run", "lint:snippets"}).
		Stdout(ctx)
}

func (s *Shopware) Phpstan(ctx context.Context) (string, error) {
	return s.
		BasicStack(ctx).
		With(WithTestEnvironment()).
		With(WithInstall()).
		WithExec([]string{"composer", "run", "framework:schema:dump"}).
		WithExec([]string{"composer", "run", "phpstan"}).
		Stdout(ctx)
}

// Executes a phpunit test suite
func (s *Shopware) Phpunit(
	ctx context.Context,
	// the name of the testsuite to run.
	// possible values are: unit, migration
	// +optional
	// +default="unit"
	testsuite string,
	// a subpath the testsuite should be limited to
	// possible values are: tests/integration/Administration, tests/integration/Core/Checkout, tests/integration/Core/Content, tests/integration/Core/Framework, tests/integration/Core/Installer, tests/integration/Core/Maintenance, tests/integration/Core/System, tests/integration/Elasticsearch, tests/integration/Storefront, .
	// +optional
	path string,
	// whether to run the testsuite with coverage
	// +optional
	// +default=false
	coverage bool,
	// whether to stop the test execution after the first failure
	// +optional
	// +default=true
	stopOnFailure bool,
) *dagger.Container {
	cmd := []string{"php", "-d", "memory_limit=-1", "vendor/bin/phpunit", "-d", "error_reporting=E_ALL", "--testsuite", testsuite}

	if stopOnFailure {
		cmd = append(cmd, "--stop-on-failure")
	}

	if coverage {
		cmd = append(cmd, "--coverage-cobertura", "coverage.xml")
	}

	if path != "" {
		cmd = append(cmd, "--", path)
	}

	c := s.
		BasicStack(ctx).
		With(WithTestEnvironment())

	if testsuite != "unit" {
		c = c.With(WithTestInstall())
	}

	if testsuite == "integration" {
		return c.
			With(WithWebServerAndExec(cmd))
	}

	return c.
		WithExec(cmd)
}

func (s *Shopware) PhpunitUnit(
	ctx context.Context,
	// whether to run the testsuite with coverage
	// +optional
	// +default=false
	coverage bool,
) (string, error) {
	return s.
		Phpunit(ctx, "unit", "", coverage, true).
		Stdout(ctx)
}

func (s *Shopware) PhpunitMigration(
	ctx context.Context,
	// whether to run the testsuite with coverage
	// +optional
	// +default=false
	coverage bool,
) (string, error) {
	return s.
		Phpunit(ctx, "migration", "", coverage, true).
		Stdout(ctx)
}

func (s *Shopware) PhpunitIntegration(
	ctx context.Context,
	// the name of the testsuite to run
	// possible values are: administration, checkout, content, framework, installer, maintenance, system, elasticsearch, storefront
	testsuite string,
	// whether to run the testsuite with coverage
	// +optional
	// +default=false
	coverage bool,
) (string, error) {
	for _, path := range integrationTestPaths {
		if strings.HasSuffix(strings.ToLower(path), testsuite) {
			return s.
				Phpunit(ctx, "integration", path, coverage, true).
				Stdout(ctx)
		}
	}

	return "", fmt.Errorf("unknown testsuite: %s", testsuite)
}

// Run all integration tests consecutively
func (s *Shopware) PhpunitIntegrationAll(ctx context.Context) (string, error) {
	var output string

	for _, path := range integrationTestPaths {
		out, err := s.
			Phpunit(ctx, "integration", path, false, true).
			Stdout(ctx)
		if err != nil {
			return output, err
		}
		output += out
	}

	return output, nil
}

// !EXPERIMENTAL - Run all integration tests in parallel
func (s *Shopware) PhpunitIntegrationAllParallel(
	ctx context.Context,
	// number of parallel tests to run
	// +optional
	// +default=2
	parallel int,
) (string, error) {
	var output string

	// TODO: Use a distinct DB for each test suite

	g, gctx := errgroup.WithContext(ctx)
	g.SetLimit(parallel) // Limit the number of parallel tests

	for _, path := range integrationTestPaths {
		path := path

		g.Go(func() error {
			out, err := s.
				Phpunit(gctx, "integration", path, false, true).
				Stdout(gctx)
			if err != nil {
				return err
			}

			output += out

			return nil
		})
	}

	return output, g.Wait()
}

func WithTestEnvironment() dagger.WithContainerFunc {
	return func(c *dagger.Container) *dagger.Container {
		return c.
			With(EnvVariables(map[string]string{
				"APP_ENV":               "test",
				"BLUE_GREEN_DEPLOYMENT": "1",
			}))
	}
}

func WithCsFixerCache() dagger.WithContainerFunc {
	return func(c *dagger.Container) *dagger.Container {
		return c.
			WithMountedCache(shopwareProjectRoot+"/var/cache/cs_fixer", dag.CacheVolume("cs_fixer"), dagger.ContainerWithMountedCacheOpts{
				Owner: shopwareUser,
			})
	}
}
