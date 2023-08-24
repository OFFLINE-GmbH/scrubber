package scrubber

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
)

// zipAction represents the action of zipping up old files.
type zipAction struct {
	action
}

// newZipAction returns a pointer to a deleteAction.
func newZipAction(dir *directory, fs Filesystem, log logger, pretend bool) *zipAction {
	return &zipAction{action{dir, fs, log, pretend}}
}

// perform zips files that are past a certain age or certain size.
func (a zipAction) perform(files []os.FileInfo, check checkFn) ([]os.FileInfo, error) {
	for i, file := range files {
		file := file
		filename := a.fs.FullPath(file, a.dir.Path)

		if check(file) {
			if a.pretend {
				a.log.Printf("[ZIP] PRETEND: Would zip file %s", filename)
				continue
			}

			a.log.Printf("[ZIP] Zipping file %s", filename)

			err := a.zip(filename)
			if err != nil {
				return files, err
			}
			err = a.fs.Remove(filename)
			if err != nil {
				a.log.Printf("[ZIP] ERROR: Failed to delete original file %s: %s", filename, err)
				continue
			}
			files = a.action.removeFile(files, i)
		} else {
			a.log.Printf("[ZIP] No action is needed for file %s", filename)
		}
	}
	return files, nil
}

// zip creates a zip file containing a single file.
func (a zipAction) zip(filePath string) error {
	zipName := filePath + ".zip"
	zipFile, err := a.fs.Create(zipName)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %v", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	info, err := a.fs.Stat(filePath)
	if err != nil {
		return nil
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf("failed to create zip header: %v", err)
	}

	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("failed to create zip writer: %v", err)
	}

	file, err := a.fs.Open(filePath)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("failed to open zip file: %v", err)
	}

	_, err = io.Copy(writer, file)

	return err
}
