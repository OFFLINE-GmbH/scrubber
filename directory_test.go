package scrubber

import (
	"os"
	"testing"
)

// setupDirectoryTest creates an in-memory filesystem and adds some files.
func setupDirectoryTest() {
	// appFS := afero.NewMemMapFs()
	//
	// appFS.MkdirAll(dir, 0755)
	// afero.WriteFile(appFS, dir+"/include.txt", []byte("file a"), 0644)
	// afero.WriteFile(appFS, dir+"/include.pdf", []byte("file a"), 0644)
	// afero.WriteFile(appFS, dir+"/exclude.exe", []byte("file b"), 0644)
	// afero.WriteFile(appFS, dir+"/exclude.zip", []byte("file b"), 0644)
	//
	// return appFS
}

// TestExclude checks if exclude rules are applied correctly.
func TestExclude(t *testing.T) {
	fs := &mockedFs{
		files: []os.FileInfo{
			mockedFileInfo{name: "include.txt"},
			mockedFileInfo{name: "include.pdf"},
			mockedFileInfo{name: "exclude.exe"},
			mockedFileInfo{name: "exclude.zip"},
		},
	}

	d := newDirectoryScanner(&directory{
		Name: "Logs", Path: testPath, Exclude: []string{"zip", "exe"},
	}, fs)

	files, err := d.getFiles()
	if err != nil {
		t.Errorf("Failed to load files: %s", err)
	}
	files = d.filterFiles(files)
	checkInclusion(files, t)
}

// TestInclude checks if exclude rules are applied correctly.
func TestInclude(t *testing.T) {
	fs := &mockedFs{
		files: []os.FileInfo{
			mockedFileInfo{name: "include.txt"},
			mockedFileInfo{name: "include.pdf"},
			mockedFileInfo{name: "exclude.exe"},
			mockedFileInfo{name: "exclude.zip"},
		},
	}

	d := newDirectoryScanner(&directory{
		Name: "Logs", Path: testPath, Include: []string{"txt", "pdf"},
	}, fs)

	files, err := d.getFiles()
	if err != nil {
		t.Errorf("Failed to load files: %s", err)
	}
	files = d.filterFiles(files)

	checkInclusion(files, t)
}

// checkInclusion is a helper function to check if the include/exclude rules were applied correctly.
func checkInclusion(files []os.FileInfo, t *testing.T) {
	gotIncludeTxt := false
	gotIncludePdf := false
	for _, f := range files {
		if f.Name() == "include.txt" {
			gotIncludeTxt = true
		}
		if f.Name() == "include.pdf" {
			gotIncludePdf = true
		}
		if f.Name() == "exclude.exe" || f.Name() == "exclude.zip" {
			t.Errorf("file \"%s\" shoud be excluded.\n", f.Name())
		}
	}
	if !gotIncludeTxt {
		t.Errorf("file 'include.txt' shoud be included.\n")
	}
	if !gotIncludePdf {
		t.Errorf("file 'include.pdf' shoud be included.\n")
	}
}
