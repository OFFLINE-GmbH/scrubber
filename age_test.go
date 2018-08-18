package scrubber

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

// TestAge tests that old files are being deleted correctly.
func TestAge(t *testing.T) {

	files := []os.FileInfo{
		mockedFileInfo{name: "dayold", modTime: time.Now().AddDate(0, 0, -1)},
		mockedFileInfo{name: "weekold", modTime: time.Now().AddDate(0, 0, -7)},
	}

	fs := &mockedFs{}

	c := StrategyConfig{Type: StrategyTypeAge, Limit: "7d", Action: ActionTypeDelete}
	d := directory{Path: testPath}

	logger := log.New(ioutil.Discard, "", 0)

	a := newDeleteAction(&d, fs, logger, false)
	s := newAgeStrategy(&c, &d, a, logger)

	_, err := s.process(files)
	if err != nil {
		t.Errorf("expected no error, got %v\n", err)
	}

	if len(fs.deleted) != 1 || fs.deleted[0] != testPath+"/weekold" {
		t.Errorf("expected only \"weekold\" to be removed got %v.\n", fs.deleted)
	}
}

// TestAgeLimitParser checks if limit strings are unmarshaled into time.Durations correctly.
func TestAgeLimitParser(t *testing.T) {
	testTimes := []struct {
		test     string
		expected time.Duration
	}{
		{"15m", 15 * time.Minute},
		{"1h", 1 * time.Hour},
		{"1d", 24 * time.Hour},
		{"1w", 168 * time.Hour},
		{"1y", 8760 * time.Hour},
		{"2y", 2 * 8760 * time.Hour},
		{"1h 15m", 1*time.Hour + 15*time.Minute},
		{"1d 1h 15m", 25*time.Hour + 15*time.Minute},
		{"1w 2d 1h 15m", 217*time.Hour + 15*time.Minute},
		{"1y 1w 2d 1h 0m", 8977 * time.Hour},
	}

	s := newAgeStrategy(nil, &directory{Path: "test"}, nil, nil)

	for _, table := range testTimes {
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

// TestInvalidAgeLimitParser tests that invalid limit strings are being rejected.
func TestInvalidAgeLimitParser(t *testing.T) {
	testTimes := []struct {
		test string
	}{
		{"-15m"},
		{"2x"},
		{""},
		{"/"},
	}

	s := newAgeStrategy(nil, &directory{Path: "test"}, nil, nil)

	for _, table := range testTimes {
		err := s.unmarshalText([]byte(table.test))
		if err == nil {
			t.Errorf("UnmarshalText(%q) should return an error", table.test)
		}
	}
}
