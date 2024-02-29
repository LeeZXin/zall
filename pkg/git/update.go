package git

import (
	"context"
)

func UpdateServerInfo(ctx context.Context, repoPath string) error {
	_, err := NewCommand("update-server-info").Run(ctx, WithDir(repoPath))
	return err
}

func RevParse(ctx context.Context, repoPath string, args ...string) error {
	_, err := NewCommand("rev-parse").AddArgs(args...).Run(ctx, WithDir(repoPath))
	return err
}
