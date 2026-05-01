package provider

import (
	"context"
	"errors"
)

type GithubProvider struct {
	Name  string
	Owner string
	Token string
}

func (p *GithubProvider) ListRepositories(ctx context.Context) ([]Repository, error) {
	return nil, errors.New("not implemented")
}

func (p *GithubProvider) AuthenticatedURL(repo Repository) (string, error) {
	return "", errors.New("not implemented")
}
