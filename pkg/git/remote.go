package git

import (
	"context"
)

const (
	DefaultRemote = "origin"
)

func AddRemote(ctx context.Context, repoPath, name, url string, tryFetch bool) error {
	cmd := NewCommand("remote", "add", name, url)
	if tryFetch {
		cmd.AddArgs("-f")
	}
	_, err := cmd.Run(ctx, WithDir(repoPath))
	return err
}

func RemoveRemote(ctx context.Context, repoPath, name string) error {
	_, err := NewCommand("remote", "rm", name).Run(ctx, WithDir(repoPath))
	return err
}
