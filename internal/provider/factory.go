package provider

import (
	"fmt"

	"github.com/wimwenigerkind/gitmirror/internal/config"
)

func New(name string, cfg config.ProviderConfig) (Provider, error) {
	switch cfg.Type {
	case "bitbucket":
		return &BitbucketProvider{Name: name, Owner: cfg.Owner, Token: cfg.Token}, nil
	case "github":
		return &GithubProvider{Name: name, Owner: cfg.Owner, Token: cfg.Token}, nil
	default:
		return nil, fmt.Errorf("provider %q: unknown type %q", name, cfg.Type)
	}
}
