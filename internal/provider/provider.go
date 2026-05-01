package provider

import "context"

type Repository struct {
	Slug string
	URL  string
}

type Provider interface {
	ListRepositories(ctx context.Context) ([]Repository, error)
	AuthenticatedURL(repo Repository) (string, error)
}
