package main

import (
	"context"
	"log"
	"path/filepath"

	"github.com/wimwenigerkind/gitmirror/internal/config"
	"github.com/wimwenigerkind/gitmirror/internal/mirror"
	"github.com/wimwenigerkind/gitmirror/internal/provider"
	"golang.org/x/sync/errgroup"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	for name, pc := range cfg.Provider {
		p, err := provider.New(name, pc)
		if err != nil {
			log.Printf("%s: %v", name, err)
			continue
		}
		repos, err := p.ListRepositories(ctx)
		if err != nil {
			log.Printf("%s: list: %v", name, err)
			continue
		}
		var g errgroup.Group
		g.SetLimit(cfg.Concurrency)
		for _, r := range repos {
			r := r
			g.Go(func() error {
				authURL, err := p.AuthenticatedURL(r)
				if err != nil {
					log.Printf("%s/%s: auth url: %v", name, r.Slug, err)
					return nil
				}
				parts := []string{cfg.Destination, name}
				if r.Project != "" {
					parts = append(parts, r.Project)
				}
				parts = append(parts, r.Slug+".git")
				destDir := filepath.Join(parts...)
				if err := mirror.Sync(ctx, authURL, destDir); err != nil {
					log.Printf("%s/%s: sync: %v", name, r.Slug, err)
				}
				return nil
			})
		}
		_ = g.Wait()
	}
}
