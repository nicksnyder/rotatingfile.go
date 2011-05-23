package rotatingfile

import (
	"os"
)

// A Reader reads from a RotatingFile.
type Reader struct {
	RotatingFile
	beginTime, endTime, seekTime, seekOffset int64
}

// NewReader returns a new Reader that reads from a RotatingFile from beginTime to endTime.
func NewReader(secondsPerFile int, format string, beginTime, endTime int64) *Reader {
	return &Reader{RotatingFile{secondsPerFile, format}, beginTime, endTime, beginTime, 0}
}

// SeekToTime sets the offset for the next Read relative to the beginning of the File that corresponds to time.
// It returns the new time and offset, or an Error if one occured.
func (r *Reader) SeekToTime(time, offset int64) (normalizedTime, normalizedOffset int64, err os.Error) {
	r.seekTime, r.seekOffset, err = r.normalize(time, offset)
	return r.seekTime, r.seekOffset, err
}

func (r *Reader) normalize(time, offset int64) (newTime, newOffset int64, err os.Error) {
	newOffset = offset
	dir := int64(1)
	if offset < 0 {
		time -= int64(r.secondsPerFile)
		dir = -1
	}

	for newTime = time; true; newTime += int64(r.secondsPerFile) * dir {
		newTime = r.time(newTime)
		if newTime < r.beginTime || newTime > r.endTime {
			// Tried to seek past end of file
			return newTime, newOffset, os.EOF
		}

		filename := r.Filename(newTime)
		info, err := os.Lstat(filename)
		if err != nil {
			// File doesn't exist, skip it
			continue
		}

		switch dir {
		case -1:
			if -newOffset > info.Size {
				// The requested offset does not exist in this file 
				newOffset += info.Size
				continue
			}

			// Convert to a positive offset
			newOffset += info.Size
		case 1:
			if newOffset >= info.Size {
				// The requested offset does not exist in this file 
				newOffset -= info.Size
				continue
			}
		}
		break
	}
	return newTime, newOffset, nil
}

// Read reads up to len(b) bytes from the RotatingFile. It returns the number of bytes read and an Error, if any.
// EOF is signaled by a zero count with err set to EOF.
func (r *Reader) Read(b []byte) (n int, err os.Error) {
	for ; len(b) > 0; r.seekTime += int64(r.secondsPerFile) {
		if r.seekTime > r.endTime {
			// Tried to read past end of file
			return n, os.EOF
		}

		filename := r.Filename(r.seekTime)
		file, err := os.Open(filename)
		if err != nil {
			// File doesn't exist, skip it
			continue
		}

		nn, err := file.ReadAt(b[n:], r.seekOffset)
		file.Close()
		r.seekOffset += int64(nn)
		n += nn

		if err == os.EOF {
			// Continue reading from beginning of next file (if one exists)
			r.seekOffset = 0
			continue
		}
		break
	}
	return
}

// ReadAt reads up to len(b) bytes from the RotatingFile starting at the offset of the File that corresponds to time.
// It returns the number of bytes read and an Error, if any. EOF is signaled by a zero count with err set to EOF.
func (r *Reader) ReadAt(b []byte, time, offset int64) (n int, err os.Error) {
	time, offset, err = r.normalize(time, offset)
	if err != nil {
		return
	}

	for ; len(b) > 0; time += int64(r.secondsPerFile) {
		if time > r.endTime {
			// Tried to read past end of file
			return n, os.EOF
		}

		filename := r.Filename(time)
		file, err := os.Open(filename)
		if err != nil {
			// File doesn't exist, skip it
			continue
		}

		nn, err := file.ReadAt(b[n:], offset)
		file.Close()
		n += nn

		if err == os.EOF {
			// Continue reading from beginning of next file (if one exists)
			offset = 0
			continue
		}
		break
	}
	return
}
