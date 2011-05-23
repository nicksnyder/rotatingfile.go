package rotatingfile

import (
	"fmt"
)

// A RotatingFile represents a logicial file whose contents may be split
// between many actual files whose names only differ by a timestamp.
type RotatingFile struct {
	secondsPerFile int
	format         string
}

func (rf *RotatingFile) time(time int64) int64 {
	return time - time%int64(rf.secondsPerFile)
}

// Filename returns the name of the File that corresponds to the time argument.
func (rf *RotatingFile) Filename(time int64) string {
	return fmt.Sprintf(rf.format, rf.time(time))
}
