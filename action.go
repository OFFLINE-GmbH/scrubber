package scrubber

import (
	"os"
)

// checkFn is the function that determines whether a file should be cleaned up or not.
type checkFn func(file os.FileInfo) bool

// action holds a reference to the current directory and a filesystem handle
type action struct {
	dir     *directory
	fs      Filesystem
	log     logger
	pretend bool
}

// performer is the interface that wraps the single method a action implementation has to provide.
type performer interface {
	perform(files []os.FileInfo, check checkFn) ([]os.FileInfo, error)
}

// actionFromConfig returns the action defined in the configuration file.
func actionFromConfig(c *StrategyConfig, dir *directory, fs Filesystem, log logger, pretend bool) performer {
	switch c.Action {
	case ActionTypeDelete:
		return newDeleteAction(dir, fs, log, pretend)
	case ActionTypeZip:
		return newZipAction(dir, fs, log, pretend)
	default:
		log.Fatalf("Unknown action type %s", c.Type)
	}
	return nil
}

// removeFile updates the in memory list of all files we're working with.
func (a action) removeFile(files []os.FileInfo, i int) []os.FileInfo {
	if len(files) > 1 {
		files = append(files[:i], files[i+1:]...)
	} else {
		files = nil
	}
	return files
}
