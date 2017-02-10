package mock

import (
	"errors"
	"os"
	"time"
)

type FileSystem struct {
	DirectoryEntries []os.FileInfo
}

func (mock FileSystem) ReadDir(dirname string) ([]os.FileInfo, error) {

	if mock.DirectoryEntries == nil {
		return nil, errors.New("Error reading directory")
	}

	return mock.DirectoryEntries, nil
}

type FileInfo struct {
	FileName string
}

func (mock FileInfo) Name() string {
	return mock.FileName
}

func (mock FileInfo) Size() int64 {
	return 0
}

func (mock FileInfo) Mode() os.FileMode {
	return 0
}

func (mock FileInfo) ModTime() time.Time {
	return time.Time{}
}

func (mock FileInfo) IsDir() bool {
	return false
}

func (mock FileInfo) Sys() interface{} {
	return nil
}
