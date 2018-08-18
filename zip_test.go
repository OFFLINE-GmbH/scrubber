package scrubber

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

// TestZip tests if the file is zipped and removed correctly.
func TestZip(t *testing.T) {
	files := []os.FileInfo{
		mockedFileInfo{name: "filename.extension", size: 20},
	}
	fs := &mockedFs{}

	c := StrategyConfig{Type: StrategyTypeSize, Limit: "10b", Action: ActionTypeZip}
	d := directory{Path: testPath}

	logger := log.New(ioutil.Discard, "", 0)

	a := newZipAction(&d, fs, logger, false)
	s := newSizeStrategy(&c, &d, a, logger)
	_, err := s.process(files)
	if err != nil {
		t.Errorf("expected no from process error, got %v\n", err)
	}

	if len(fs.created) != 1 || fs.created[0] != testPath+"/filename.extension.zip" {
		t.Errorf("expected only \"/logs/filename.extension.zip\" to exist got %v.\n", fs.created)
	}
}
