package provider

import (
	"context"
	"errors"
)

type BitbucketProvider struct {
	Name  string
	Owner string
	Token string
}

func (p *BitbucketProvider) ListRepositories(ctx context.Context) ([]Repository, error) {
	return nil, errors.New("not implemented")
}

func (p *BitbucketProvider) AuthenticatedURL(repo Repository) (string, error) {
	return "", errors.New("not implemented")
}
