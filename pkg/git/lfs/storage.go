package lfs

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zall/util"
	"io"
	"os"
	"path/filepath"
)

// Object represents the object on the storage
type Object interface {
	io.ReadCloser
	io.Seeker
	Stat() (os.FileInfo, error)
}

type Storage interface {
	StoreDir() string
	Open(context.Context, string) (Object, error)
	Save(context.Context, string, io.Reader) (int64, error)
	Stat(context.Context, string) (os.FileInfo, error)
	Exists(context.Context, string) (bool, error)
	Delete(context.Context, string) error
	IterateObjects(context.Context, string, func(path string, obj Object) error) error
}

// localStorage represents a local files storage
type localStorage struct {
	storeDir string
	tmpdir   string
}

func (l *localStorage) StoreDir() string {
	return l.storeDir
}

// Open a file
func (l *localStorage) Open(_ context.Context, path string) (Object, error) {
	return os.Open(filepath.Join(l.storeDir, path))
}

// Save a file
func (l *localStorage) Save(ctx context.Context, path string, r io.Reader) (int64, error) {
	p := filepath.Join(l.storeDir, path)
	if err := os.MkdirAll(filepath.Dir(p), os.ModePerm); err != nil {
		return 0, err
	}
	// Create a temporary file to save to
	if err := os.MkdirAll(l.tmpdir, os.ModePerm); err != nil {
		return 0, err
	}
	tmp, err := os.CreateTemp(l.tmpdir, "upload-*")
	if err != nil {
		return 0, err
	}
	tmpRemoved := false
	defer func() {
		if !tmpRemoved {
			_ = util.RemoveAll(tmp.Name())
		}
	}()
	n, err := io.Copy(tmp, r)
	if err != nil {
		return 0, err
	}
	if err := tmp.Close(); err != nil {
		return 0, err
	}
	if err := util.Rename(tmp.Name(), p); err != nil {
		return 0, err
	}
	// Golang's tmp file (os.CreateTemp) always have 0o600 mode, so we need to change the file to follow the umask (as what Create/MkDir does)
	// but we don't want to make these files executable - so ensure that we mask out the executable bits
	if err := util.ApplyUmask(p, os.ModePerm&0o666); err != nil {
		return 0, err
	}
	tmpRemoved = true
	return n, nil
}

// Stat returns the info of the file
func (l *localStorage) Stat(_ context.Context, path string) (os.FileInfo, error) {
	return os.Stat(filepath.Join(l.storeDir, path))
}

// Exists returns the info of the file
func (l *localStorage) Exists(ctx context.Context, path string) (bool, error) {
	_, err := l.Stat(ctx, path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Delete delete a file
func (l *localStorage) Delete(_ context.Context, path string) error {
	return util.RemoveAll(filepath.Join(l.storeDir, path))
}

// IterateObjects iterates across the objects in the local storage
func (l *localStorage) IterateObjects(ctx context.Context, dirName string, fn func(path string, obj Object) error) error {
	dir := filepath.Join(l.storeDir, dirName)
	return filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if path == l.storeDir {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		relPath, err := filepath.Rel(l.storeDir, path)
		if err != nil {
			return err
		}
		obj, err := os.Open(path)
		if err != nil {
			return err
		}
		defer obj.Close()
		return fn(relPath, obj)
	})
}

func NewLocalStorage(storeDir, tmpDir string) (Storage, error) {
	if !filepath.IsAbs(storeDir) {
		return nil, fmt.Errorf("%s is not absolute path", storeDir)
	}
	if !filepath.IsAbs(tmpDir) {
		return nil, fmt.Errorf("%s is not absolute path", tmpDir)
	}
	return &localStorage{
		storeDir: storeDir,
		tmpdir:   tmpDir,
	}, nil
}
