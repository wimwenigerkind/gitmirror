package main

import (
	"context"
	"fmt"
	"log"

	"github.com/wimwenigerkind/gitmirror/internal/config"
	"github.com/wimwenigerkind/gitmirror/internal/provider"
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
			log.Fatal(err)
		}
		repos, err := p.ListRepositories(ctx)
		if err != nil {
			log.Printf("%s: %v", name, err)
			continue
		}
		for _, r := range repos {
			fmt.Println(name, r.Slug)
		}
	}
}
