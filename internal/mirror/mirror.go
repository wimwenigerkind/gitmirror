package mirror

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func Sync(ctx context.Context, authURL, destDir string) error {
	if isBareRepo(destDir) {
		return update(ctx, authURL, destDir)
	}
	return clone(ctx, authURL, destDir)
}

func clone(ctx context.Context, authURL, destDir string) error {
	if err := os.MkdirAll(filepath.Dir(destDir), 0o755); err != nil {
		return fmt.Errorf("mirror: create parent dir: %w", err)
	}
	cmd := exec.CommandContext(ctx, "git", "clone", "--mirror", authURL, destDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("mirror: clone failed: %w", err)
	}
	return nil
}

func update(ctx context.Context, authURL, destDir string) error {
	setURL := exec.CommandContext(ctx, "git", "-C", destDir, "remote", "set-url", "origin", authURL)
	setURL.Stderr = os.Stderr
	if err := setURL.Run(); err != nil {
		return fmt.Errorf("mirror: set-url failed: %w", err)
	}
	fetch := exec.CommandContext(ctx, "git", "-C", destDir, "fetch", "--prune", "origin")
	fetch.Stdout = os.Stdout
	fetch.Stderr = os.Stderr
	if err := fetch.Run(); err != nil {
		return fmt.Errorf("mirror: fetch failed: %w", err)
	}
	return nil
}

func isBareRepo(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, "HEAD"))
	return err == nil
}
