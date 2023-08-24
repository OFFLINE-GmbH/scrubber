// Package scrubber provides an easy way to clean up old files in a directory.
package scrubber

import (
	"fmt"
	"os"
	"path/filepath"
)

// Scrubber holds the configuration and a filesystem handle.
type Scrubber struct {
	config  *TomlConfig
	fs      Filesystem
	log     logger
	pretend bool
}

// TomlConfig holds the complete structure of the scrubber config file.
type TomlConfig struct {
	Title       string
	Directories []directory `toml:"directory"`
}

// Strategy represents an action to take with files.
type Strategy struct {
	c      *StrategyConfig
	dir    *directory
	action performer
	log    logger
}

// StrategyConfig holds all specified strategies for a single Directory.
type StrategyConfig struct {
	Type   StrategyType
	Action StrategyAction
	Limit  string
}

// StrategyType defines how to decide what files should be cleaned up.
type StrategyType string

const (
	// StrategyTypeAge makes files past a certain age to be deleted.
	StrategyTypeAge StrategyType = "age"
	// StrategyTypeSize makes files past a certain size to be deleted.
	StrategyTypeSize StrategyType = "size"
)

// StrategyAction represents the action that should be taken for matching files.
type StrategyAction string

const (
	// ActionTypeDelete is used to delete old files.
	ActionTypeDelete StrategyAction = "delete"
	// ActionTypeZip is used to zip old files.
	ActionTypeZip StrategyAction = "zip"
)

// processor is the interface that wraps the single method a strategy implementation has to provide.
type processor interface {
	process(files []os.FileInfo) ([]os.FileInfo, error)
}

// logger defines the the minimal logging functionality we expect.
type logger interface {
	Printf(string, ...interface{})
	Fatalf(format string, v ...interface{})
}

// New returns a new instance of Scrubber.
func New(c *TomlConfig, fs Filesystem, log logger, pretend bool) *Scrubber {
	return &Scrubber{
		config:  c,
		fs:      fs,
		log:     log,
		pretend: pretend,
	}
}

// Scrub performs the actual cleanup.
func (s Scrubber) Scrub() ([]os.FileInfo, error) {
	var files []os.FileInfo
	for _, configDir := range s.config.Directories {

		expandedDirs, err := s.expandDirs(configDir.Path)
		if err != nil {
			s.log.Printf("[ERROR] Failed to expand path %s: %s", configDir.Path, err)
			continue
		}

		if len(expandedDirs) < 1 {
			s.log.Printf("Found no files to process. Skipping %s", configDir.Path)
			continue
		}

		for _, expandedDir := range expandedDirs {
			dir := configDir.WithPath(expandedDir)

			s.log.Printf("Scanning for files in %s...", dir.Path)

			scanner := newDirectoryScanner(&dir, s.fs)
			files, err := scanner.getFiles()
			if err != nil {
				s.log.Printf("[ERROR] Failed to load files in directory %s...: %s", dir.Path, err)
				continue
			}

			files = scanner.filterFiles(files)
			files = ApplyKeepLatest(files, dir.KeepLatest)

			if len(files) < 1 {
				s.log.Printf("Found no files to process. Skipping %s", dir.Path)
				continue
			}

			s.log.Printf("Found %d files to process", len(files))

			for _, strategy := range dir.Strategies {
				s, err := strategyFromConfig(&strategy, &dir, s.fs, s.log, s.pretend)
				if err != nil {
					return nil, err
				}
				_, err = s.process(files)
				if err != nil {
					return nil, fmt.Errorf("error while processing files: %s", err)
				}
			}
		}

	}
	return files, nil
}

// expandDirs expands a Glob pattern and returns all directories.
func (s Scrubber) expandDirs(path string) ([]string, error) {
	expandedPaths, err := filepath.Glob(path)
	if err != nil {
		return nil, err
	}

	var dirs []string
	for _, expandedPath := range expandedPaths {
		fileInfo, err := s.fs.Stat(expandedPath)
		if err != nil {
			return nil, err
		}

		if fileInfo.IsDir() {
			dirs = append(dirs, expandedPath)
		}
	}

	return dirs, nil
}

// strategyFromConfig returns the strategy defined in the configuration file.
func strategyFromConfig(c *StrategyConfig, dir *directory, fs Filesystem, log logger, pretend bool) (processor,
	error) {
	action := actionFromConfig(c, dir, fs, log, pretend)
	switch c.Type {
	case StrategyTypeAge:
		return newAgeStrategy(c, dir, action, log), nil
	case StrategyTypeSize:
		return newSizeStrategy(c, dir, action, log), nil
	}
	return nil, fmt.Errorf("unknown strategy type: %s", c.Type)
}
