package io

import (
	"io/ioutil"
	"os"
)

type File interface {
	Read(p []byte) (n int, err error)
}

type fileSystem interface {
	Open(name string) (File, error)
	ReadDir(dirname string) ([]os.FileInfo, error)
}

type realFileSystem struct {
}

func (realFileSystem) Open(name string) (File, error) {
	return os.Open(name)
}

func (realFileSystem) ReadDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dirname)
}

var FileSystem fileSystem = realFileSystem{}
