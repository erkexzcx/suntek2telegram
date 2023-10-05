package ftpserver

import (
	"io"
	"io/fs"
	"os"
	"time"

	"goftp.io/server/v2"
)

func (perm *MyPerm) GetOwner(s string) (string, error) {
	return "", nil
}
func (perm *MyPerm) GetGroup(s string) (string, error) {
	return "", nil
}
func (perm *MyPerm) GetMode(s string) (os.FileMode, error) {
	return 0644, nil
}

func (perm *MyPerm) ChOwner(s string, asdf string) error {
	return nil
}
func (perm *MyPerm) ChGroup(s string, asdf string) error {
	return nil
}
func (perm *MyPerm) ChMode(s string, asdf os.FileMode) error {
	return nil
}

type DummyFileInfo struct{}

func (dfi *DummyFileInfo) Name() string       { return "dummy" }
func (dfi *DummyFileInfo) Size() int64        { return 0 }
func (dfi *DummyFileInfo) Mode() fs.FileMode  { return 0 }
func (dfi *DummyFileInfo) ModTime() time.Time { return time.Time{} }
func (dfi *DummyFileInfo) IsDir() bool        { return false }
func (dfi *DummyFileInfo) Sys() interface{}   { return nil }

type DummyDirInfo struct{}

func (ddi *DummyDirInfo) Name() string       { return "dummy" }
func (ddi *DummyDirInfo) Size() int64        { return 0 }
func (ddi *DummyDirInfo) Mode() fs.FileMode  { return fs.ModeDir }
func (ddi *DummyDirInfo) ModTime() time.Time { return time.Time{} }
func (ddi *DummyDirInfo) IsDir() bool        { return true }
func (ddi *DummyDirInfo) Sys() interface{}   { return nil }

var currentDir = "/"

func (driver *MyDriver) Stat(ctx *server.Context, path string) (fs.FileInfo, error) {
	if path == currentDir {
		return &DummyDirInfo{}, nil
	}
	return &DummyFileInfo{}, nil
}
func (driver *MyDriver) ChangeDir(ctx *server.Context, path string) error {
	return nil
}
func (driver *MyDriver) ListDir(ctx *server.Context, path string, callback func(fs.FileInfo) error) error {
	return nil
}
func (driver *MyDriver) DeleteDir(ctx *server.Context, path string) error {
	return nil
}
func (driver *MyDriver) DeleteFile(ctx *server.Context, path string) error {
	return nil
}
func (driver *MyDriver) Rename(ctx *server.Context, fromPath string, toPath string) error {
	return nil
}
func (driver *MyDriver) MakeDir(ctx *server.Context, path string) error {
	return nil
}

type DummyReadCloser struct{}

func (drc *DummyReadCloser) Read(p []byte) (n int, err error) {
	return 0, io.EOF // return EOF to signal end of file
}

func (drc *DummyReadCloser) Close() error {
	return nil
}

func (driver *MyDriver) GetFile(ctx *server.Context, path string, offset int64) (int64, io.ReadCloser, error) {
	return 0, &DummyReadCloser{}, nil
}
