package main

import (
	"dagger/shopware/internal/dagger"
)

func WithInstall() dagger.WithContainerFunc {
	return func(c *dagger.Container) *dagger.Container {
		return c.
			WithExec([]string{"bin/console", "system:install", "--basic-setup", "--force"})
	}
}

func WithFullInstall() dagger.WithContainerFunc {
	return func(c *dagger.Container) *dagger.Container {
		return c.
			WithExec([]string{"composer", "setup"})
	}
}

func WithTestInstall() dagger.WithContainerFunc {
	return func(c *dagger.Container) *dagger.Container {
		return c.
			WithExec([]string{"php", "src/Core/TestBootstrap.php"})
	}
}
