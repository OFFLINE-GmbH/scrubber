package scrubber

import (
	"fmt"
	"os"

	"github.com/c2h5oh/datasize"
)

// ageStrategy represents the action of deleting files that are greater than a certain file size.
type sizeStrategy struct {
	Strategy
	limit int64
}

// newSizeStrategy returns a new sizeStrategy.
func newSizeStrategy(c *StrategyConfig, dir *directory, action performer, log logger) *sizeStrategy {
	return &sizeStrategy{Strategy{c, dir, action, log}, 0}
}

// process deletes all files greater than a certain size.
func (s sizeStrategy) process(files []os.FileInfo) ([]os.FileInfo, error) {
	err := s.unmarshalText([]byte(s.c.Limit))
	if err != nil {
		return nil, err
	}
	if s.limit <= 0 {
		return nil, fmt.Errorf("file size limit has to be greater than 0")
	}

	files, err = s.action.perform(files, func(file os.FileInfo) bool {
		return file.Size() > s.limit
	})

	return files, err
}

// parseLimit turns a string representation of a filesize into bytes
func (s *sizeStrategy) unmarshalText(text []byte) error {
	if len(text) < 1 {
		return fmt.Errorf("limit cannot be an empty string")
	}

	var v datasize.ByteSize
	err := v.UnmarshalText(text)
	if err != nil {
		return fmt.Errorf("invalid size definition")
	}

	s.limit = int64(v)

	return nil
}
