package main

import (
	"dagger/shopware/internal/dagger"
)

type Shopware struct {
	Source *dagger.Directory
}

func New(
	// The shopware source code
	// +optional
	Source *dagger.Directory,
) *Shopware {
	if Source == nil {
		Source = dag.Git(
			"https://github.com/shopware/shopware.git",
			dagger.GitOpts{KeepGitDir: true},
		).
			Branch("trunk").
			Tree()
	}

	return &Shopware{
		Source: Source,
	}
}