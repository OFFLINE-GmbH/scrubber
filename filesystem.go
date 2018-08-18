package scrubber

import (
	"fmt"
	"os"
	"path"
)

// Filesystem represents the minimal fs implementation we expect.
type Filesystem interface {
	Name(file os.FileInfo) string
	FullPath(file os.FileInfo, dir string) string
	Remove(path string) error
	Open(name string) (*os.File, error)
	Create(name string) (*os.File, error)
	Stat(name string) (os.FileInfo, error)
	ListFiles(path string) ([]os.FileInfo, error)
	Ext(file os.FileInfo) string
}

// OSFilesystem proxies calls to the underlying os and file library calls.
type OSFilesystem struct {
}

// Name returns the name of a file.
func (fs OSFilesystem) Name(file os.FileInfo) string {
	return file.Name()
}

// FullPath combines a file's name and it's path to a full path string.
func (fs OSFilesystem) FullPath(file os.FileInfo, dir string) string {
	return dir + "/" + file.Name()
}

// Remove deletes a file from the filesystem.
func (fs OSFilesystem) Remove(path string) error {
	return os.Remove(path)
}

// Open reads a file from the filesystem.
func (fs OSFilesystem) Open(name string) (*os.File, error) {
	return os.Open(name)
}

// Create creates a file on the filesystem.
func (fs OSFilesystem) Create(name string) (*os.File, error) {
	return os.Create(name)
}

// Stat returns information to a specific file.
func (fs OSFilesystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

// Ext returns a file's extension.
func (fs OSFilesystem) Ext(file os.FileInfo) string {
	return path.Ext(file.Name())
}

// ListFiles returns a os.FileInfo for every file in a directory.
func (fs OSFilesystem) ListFiles(path string) ([]os.FileInfo, error) {
	d, err := os.Open(path)
	defer d.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %s", path, err)
	}
	files, err := d.Readdir(-1)
	if err != nil {
		return nil, fmt.Errorf("failed to read files from directory %s: %s", path, err)
	}
	return files, nil
}
