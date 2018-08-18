package scrubber

import (
	"log"
	"os"
	"testing"
)

// TestSize tests that big files are being deleted correctly.
func TestSize(t *testing.T) {

	files := []os.FileInfo{
		mockedFileInfo{name: "20bytes", size: 20},
		mockedFileInfo{name: "10bytes", size: 10},
	}

	fs := &mockedFs{}

	c := StrategyConfig{Type: StrategyTypeSize, Limit: "10b", Action: ActionTypeDelete}
	d := directory{Path: testPath}

	logger := log.New(os.Stdout, "", 0)

	a := newDeleteAction(&d, fs, logger, false)
	s := newSizeStrategy(&c, &d, a, logger)
	_, err := s.process(files)
	if err != nil {
		t.Errorf("expected no error, got %v\n", err)
	}

	if len(fs.deleted) != 1 || fs.deleted[0] != testPath+"/20bytes" {
		t.Errorf("expected only \"20bytes\" to be removed got %v.\n", fs.deleted)
	}
}

// TestSizeLimitParser checks if limit strings are unmarshaled into int64 correctly.
func TestSizeLimitParser(t *testing.T) {
	testSizes := []struct {
		test     string
		expected int64
	}{
		{"10MB", int64(10 * 1024 * 1024)},
		{"1G", int64(1024 * 1024 * 1024)},
		{"10b", int64(10)},
	}

	s := newSizeStrategy(nil, &directory{Path: "test"}, nil, log.New(os.Stdout, "", 0))

	for _, table := range testSizes {
		err := s.unmarshalText([]byte(table.test))
		if err != nil {
			t.Errorf("UnmarshalText(%q) = %q. got unexpected error %q",
				table.test,
				s.limit,
				err,
			)
		}
		if s.limit != table.expected {
			t.Errorf("UnmarshalText(%q) = %q. got %q with error %q, expected %q",
				table.test,
				s.limit,
				err,
				s.limit,
				table.expected,
			)
		}
	}
}

// TestInvalidSizeLimitParser tests that invalid limit strings are being rejected.
func TestInvalidSizeLimitParser(t *testing.T) {
	testSizes := []struct {
		test string
	}{
		{"-10M"},
		{"2x"},
		{""},
		{"/"},
	}

	s := newSizeStrategy(nil, &directory{Path: "test"}, nil, log.New(os.Stdout, "", 0))

	for _, table := range testSizes {
		err := s.unmarshalText([]byte(table.test))
		if err == nil {
			t.Errorf("UnmarshalText(%q) should return an error", table.test)
		}
	}
}
