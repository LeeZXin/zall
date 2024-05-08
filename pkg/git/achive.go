package git

import (
	"context"
	"github.com/pingcap/errors"
	"io"
)

// ArchiveType archive types
type ArchiveType string

const (
	// ZIP zip archive type
	ZIP ArchiveType = "zip"
	// TARGZ tar gz archive type
	TARGZ ArchiveType = "tar.gz"
)

func (t ArchiveType) HttpContentType() string {
	switch t {
	case ZIP:
		return "application/zip"
	case TARGZ:
		return "application/x-tar"
	default:
		return ""
	}
}

func (t ArchiveType) String() string {
	return string(t)
}

func (t ArchiveType) IsValid() bool {
	switch t {
	case ZIP, TARGZ:
		return true
	default:
		return false
	}
}

func CreateArchive(ctx context.Context, repoPath, commitId string, t ArchiveType, writer io.Writer) error {
	if writer == nil {
		return errors.New("empty writer")
	}
	return NewCommand("archive", "--format="+t.String()).
		AddDynamicArgs(commitId).
		RunWithStdout(ctx, writer, WithDir(repoPath))
}
