package scrubber

import (
	"os"
)

// deleteAction represents the action of deleting old or big files.
type deleteAction struct {
	action
}

// newDeleteAction returns a pointer to a deleteAction.
func newDeleteAction(dir *directory, fs Filesystem, log logger, pretend bool) *deleteAction {
	return &deleteAction{action{dir, fs, log, pretend}}
}

// perform deletes files that are past a certain age or certain size.
func (a deleteAction) perform(files []os.FileInfo, check checkFn) ([]os.FileInfo, error) {
	for i, file := range files {
		filename := a.fs.FullPath(file, a.dir.Path)
		if check(file) {
			if a.pretend {
				a.log.Printf("[Delete] PRETEND: Would delete file %s", filename)
				continue
			}

			a.log.Printf("[Delete] Deleting file %s", filename)

			err := a.fs.Remove(filename)

			if err != nil {
				a.log.Printf("[Delete] ERROR: Failed to delete file %s: %s", filename, err)
				continue
			}
			files = a.action.removeFile(files, i)
		} else {
			a.log.Printf("[Delete] No action is needed for file %s", filename)
		}
	}
	return files, nil
}
