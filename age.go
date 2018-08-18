package scrubber

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// durationUnits maps the string representation of a unit to a time.Duration.
var durationUnits = map[string]time.Duration{
	"m": time.Minute,
	"h": time.Hour,
	"d": 24 * time.Hour,
	"w": 7 * 24 * time.Hour,
	"y": 8760 * time.Hour,
}

// ageStrategy represents the action of deleting files past a certain age.
type ageStrategy struct {
	Strategy
	limit time.Duration
}

// newAgeStrategy returns a new ageStrategy.
func newAgeStrategy(c *StrategyConfig, dir *directory, action performer, log logger) *ageStrategy {
	return &ageStrategy{Strategy{c, dir, action, log}, 0}
}

// process deletes all files past a certain age.
func (s ageStrategy) process(files []os.FileInfo) ([]os.FileInfo, error) {
	err := s.unmarshalText([]byte(s.c.Limit))
	if err != nil {
		return nil, err
	}

	deadline := time.Now().Add(-1 * s.limit)

	files, err = s.action.perform(files, func(file os.FileInfo) bool {
		return file.ModTime().Before(deadline)
	})

	return files, err
}

// unmarshalText turns a string representation of a duration into a time.Duration
func (s *ageStrategy) unmarshalText(text []byte) error {
	var total time.Duration
	var limit = string(text)

	if limit == "" {
		return fmt.Errorf("limit cannot be an empty string")
	}

	for _, part := range strings.Fields(limit) {
		quantifier, err := strconv.Atoi(part[:len(part)-1])
		if err != nil || quantifier < 0 {
			return fmt.Errorf("invalid limit quantifier %q passed (%s)", quantifier, err)
		}

		duration, ok := durationUnits[part[len(part)-1:]]
		if !ok {
			return fmt.Errorf("unknown limit unit %s passed to AgeStrategy", part[:len(part)-1])
		}

		total = total + (time.Duration(quantifier) * duration)
	}

	s.limit = total

	return nil
}
