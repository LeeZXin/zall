package git

import (
	"context"
	"io"
)

func UploadPack(ctx context.Context, repoPath string, input io.Reader, output io.Writer, env []string) error {
	return NewCommand("upload-pack", "--stateless-rpc", repoPath).RunWithStdout(ctx,
		output,
		WithDir(repoPath),
		WithStdin(input),
		WithEnv(env),
	)
}

func ReceivePack(ctx context.Context, repoPath string, input io.Reader, output io.Writer, env []string) error {
	return NewCommand("receive-pack", "--stateless-rpc", repoPath).RunWithStdout(ctx,
		output,
		WithDir(repoPath),
		WithStdin(input),
		WithEnv(env),
	)
}

func InfoRefs(ctx context.Context, repoPath, service string, env []string) ([]byte, error) {
	result, err := NewCommand(service, "--stateless-rpc", "--advertise-refs", ".").
		Run(ctx, WithDir(repoPath), WithEnv(env))
	if err != nil {
		return nil, err
	}
	return result.ReadAsBytes(), nil
}
