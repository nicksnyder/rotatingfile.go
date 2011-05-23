package rotatingfile

import (
	"os"
)

type Reader struct {
	RotatingFile
	begin, end, seekTime, seekOffset int64
}

func NewReader(secondsPerFile int, format string, begin, end int64) *Reader {
	return &Reader{RotatingFile{secondsPerFile, format}, begin, end, begin, 0}
}

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

	for i := int64(0); true; i += dir {
		newTime = r.time(time, i)
		if newTime < r.begin || newTime > r.end {
			// Tried to seek past end of file
			return newTime, newOffset, os.EOF
		}

		filename := r.filename(time, i)
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

func (r *Reader) Read(b []byte) (n int, err os.Error) {
	for ; len(b) > 0; r.seekTime = r.time(r.seekTime, 1) {
		if r.seekTime > r.end {
			// Tried to read past end of file
			return n, os.EOF
		}

		filename := r.filename(r.seekTime, 0)
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

func (r *Reader) ReadAt(b []byte, time, offset int64) (n int, err os.Error) {
	time, offset, err = r.normalize(time, offset)
	if err != nil {
		return
	}

	for ; len(b) > 0; time = r.time(time, 1) {
		if time > r.end {
			// Tried to read past end of file
			return n, os.EOF
		}

		filename := r.filename(time, 0)
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
