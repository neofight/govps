package mock

import (
	"errors"
	stdio "io"
	"os"
	"strings"
	"time"

	"github.com/neofight/govps/io"
)

type File struct {
	reader stdio.Reader
}

func (mock File) Read(p []byte) (n int, err error) {
	return mock.reader.Read(p)
}

type FileSystem struct {
	DirectoryEntries []os.FileInfo
	Files            map[string]string
}

func NewFileSystem() *FileSystem {
	return &FileSystem{Files: make(map[string]string)}
}

func (mock FileSystem) Open(name string) (io.File, error) {

	contents, ok := mock.Files[name]

	if !ok {
		return File{}, errors.New("Error reading file")
	}

	return File{strings.NewReader(contents)}, nil
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
