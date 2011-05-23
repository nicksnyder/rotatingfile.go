package rotatingfile

import (
	"fmt"
)

// A RotatingFile represents a logicial file whose contents may be split
// between many actual files and whose name only differs by a timestamp
type RotatingFile struct {
	secondsPerFile int
	format         string
}

func (rf *RotatingFile) time(time, index int64) int64 {
	return time - time%int64(rf.secondsPerFile) + index*int64(secondsPerFile)
}

func (rf *RotatingFile) filename(time, index int64) string {
	return fmt.Sprintf(rf.format, rf.time(time, index))
}
