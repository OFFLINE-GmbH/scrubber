package scrubber

import (
	"os"
	"strings"
	"time"
)

// testPath is a "dummy" path used for all files during testing.
var testPath = "/logs"

// mockedFs implements the Filesystem interface for testing.
type mockedFs struct {
	OSFilesystem
	files   []os.FileInfo
	deleted []string
	created []string
}

// Remove marks a file as removed on the mocked filesystem.
func (fs *mockedFs) Remove(path string) error {
	fs.deleted = append(fs.deleted, path)
	return nil
}

// Create marks a file as created on the mocked filesystem.
func (fs *mockedFs) Create(name string) (*os.File, error) {
	fs.created = append(fs.created, name)
	return nil, nil
}

// ListFiles returns all mocked files.
func (fs mockedFs) ListFiles(path string) ([]os.FileInfo, error) {
	return fs.files, nil
}

// Ext returns the file extension for a certain file.
func (fs mockedFs) Ext(file os.FileInfo) string {
	return "." + strings.Split(file.Name(), ".")[1]
}

// mockedFileInfo is a os.FileInfo implementation used for testing.
type mockedFileInfo struct {
	os.FileInfo
	modTime time.Time
	size    int64
	name    string
}

// Name returns the mocked name of a file.
func (m mockedFileInfo) Name() string { return m.name }

// Size returns the mocked size of a file.
func (m mockedFileInfo) Size() int64 { return m.size }

// ModTime returns the mocked modification time of a file.
func (m mockedFileInfo) ModTime() time.Time { return m.modTime }

// Name returns the mocked mode of a file.
func (m mockedFileInfo) Mode() os.FileMode { return 0 }
