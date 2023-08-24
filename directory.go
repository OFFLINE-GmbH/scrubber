package scrubber

import (
	"os"
	"slices"
)

// directory holds the cleanup information for a single path in the filesystem.
type directory struct {
	Name       string
	Path       string
	Include    []string
	Exclude    []string
	Strategies []StrategyConfig `toml:"strategy"`
	KeepLatest int
}

// WithPath returns a copy of the struct with the Path field set to dir.
func (d directory) WithPath(dir string) directory {
	return directory{
		Name:       d.Name,
		Path:       dir,
		Include:    d.Include,
		Exclude:    d.Exclude,
		Strategies: d.Strategies,
	}
}

// directoryScanner is used to scan a directory for files.
type directoryScanner struct {
	dir *directory
	fs  Filesystem
}

// newDirectoryScanner returns a pointer to a directoryScanner.
func newDirectoryScanner(dir *directory, fs Filesystem) *directoryScanner {
	return &directoryScanner{
		dir,
		fs,
	}
}

// getFiles returns all files in the cleanup directory.
func (s directoryScanner) getFiles() ([]os.FileInfo, error) {
	return s.fs.ListFiles(s.dir.Path)
}

// filterFiles applies the include and exclude rules to all files in the cleanup directory.
func (s directoryScanner) filterFiles(files []os.FileInfo) []os.FileInfo {
	var filtered []os.FileInfo
	for _, file := range files {
		if !file.Mode().IsRegular() {
			continue
		}

		fileExt := s.fs.Ext(file)

		var include bool
		if s.dir.Include != nil {
			include = includesExtension(s.dir.Include, fileExt)
		} else {
			include = !includesExtension(s.dir.Exclude, fileExt)
		}

		if include {
			filtered = append(filtered, file)
		}

	}
	return filtered
}

// includesExtension checks if a certain extension is an element of extensions.
func includesExtension(extensions []string, fileExt string) bool {
	for _, ext := range extensions {
		if "."+ext == fileExt {
			return true
		}
	}
	return false
}

// ApplyKeepLatest applies the keep latest rule to a slice of files.
func ApplyKeepLatest(files []os.FileInfo, latest int) []os.FileInfo {
	if latest < 1 {
		return files
	}

	slices.SortFunc(files, func(i, j os.FileInfo) int {
		if i.ModTime().After(j.ModTime()) {
			return 0
		}
		return 1
	})

	if len(files) > latest {
		return files[latest:]
	}

	return files
}
