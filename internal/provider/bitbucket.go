package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type BitbucketProvider struct {
	Name  string
	Owner string
	Token string
}

var bitbucketClient = &http.Client{Timeout: 30 * time.Second}

type bitbucketListResponse struct {
	Values []bitbucketRepo `json:"values"`
	Next   string          `json:"next"`
}

type bitbucketRepo struct {
	Slug    string `json:"slug"`
	Project struct {
		Key string `json:"key"`
	} `json:"project"`
	Links struct {
		Clone []struct {
			Name string `json:"name"`
			Href string `json:"href"`
		} `json:"clone"`
	} `json:"links"`
}

func (p *BitbucketProvider) ListRepositories(ctx context.Context) ([]Repository, error) {
	var repos []Repository
	url := "https://api.bitbucket.org/2.0/repositories/" + p.Owner

	for url != "" {
		page, err := p.fetchPage(ctx, url)
		if err != nil {
			return nil, err
		}
		for _, r := range page.Values {
			repos = append(repos, Repository{
				Slug:    r.Slug,
				URL:     httpsCloneURL(r),
				Project: r.Project.Key,
			})
		}
		url = page.Next
	}
	return repos, nil
}

func (p *BitbucketProvider) fetchPage(ctx context.Context, url string) (*bitbucketListResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if p.Token != "" {
		req.Header.Set("Authorization", "Bearer "+p.Token)
	}
	req.Header.Set("Accept", "application/json")

	res, err := bitbucketClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("bitbucket: unexpected status %d", res.StatusCode)
	}

	var page bitbucketListResponse
	if err := json.NewDecoder(res.Body).Decode(&page); err != nil {
		return nil, fmt.Errorf("bitbucket: decode: %w", err)
	}
	return &page, nil
}

func httpsCloneURL(repo bitbucketRepo) string {
	for _, c := range repo.Links.Clone {
		if c.Name == "https" {
			return c.Href
		}
	}
	return ""
}

func (p *BitbucketProvider) AuthenticatedURL(repo Repository) (string, error) {
	if repo.URL == "" {
		return "", fmt.Errorf("bitbucket: repository %q has no clone URL", repo.Slug)
	}
	u, err := url.Parse(repo.URL)
	if err != nil {
		return "", fmt.Errorf("bitbucket: parse clone URL: %w", err)
	}
	if p.Token != "" {
		u.User = url.UserPassword("x-token-auth", p.Token)
	}
	return u.String(), nil
}
