package main

import (
	"dagger/shopware/internal/dagger"
)

func WithElasticsearch() dagger.WithContainerFunc {
	elasticsearch := dag.Container().From("opensearchproject/opensearch:1").
		With(EnvVariables(map[string]string{
			"discovery.type":            "single-node",
			"plugins.security.disabled": "true",
		})).
		WithExposedPort(9200).
		AsService()

	return func(c *dagger.Container) *dagger.Container {
		return c.
			WithServiceBinding("elasticsearch", elasticsearch).
			With(EnvVariables(map[string]string{
				"OPENSEARCH_URL": "elasticsearch:9200",
			}))

	}
}
