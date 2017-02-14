package io

import (
	"io/ioutil"
	"os"
)

type fileSystem interface {
	ReadDir(dirname string) ([]os.FileInfo, error)
}

type realFileSystem struct {
}

func (realFileSystem) ReadDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dirname)
}

var FileSystem fileSystem = realFileSystem{}
