package rotatingfile

import (
	"os"
	"time"
)

// A Writer writes to a RotatingFile.
type Writer struct {
	RotatingFile
	file *os.File // The current file being written to
}

// NewWriter returns a new Writer that writes to a new File every secondsPerFile seconds.
// New Files are named according to the format string (e.g. "error.%v.log").
func NewWriter(secondsPerFile int, format string) *Writer {
	return &Writer{RotatingFile{secondsPerFile, format}, nil}
}

func (w *Writer) open(time int64) (err os.Error) {
	filename := w.Filename(time)

	// Already have the file open
	if w.file != nil && w.file.Name() == filename {
		return nil
	}

	// Close the old log file
	if w.file != nil {
		w.file.Close()
	}

	// Rotate to new file
	w.file, err = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	return err
}

// Write is equivalent to WriteAtTime(p, time.Seconds()).
func (w *Writer) Write(p []byte) (n int, err os.Error) {
	return w.WriteAtTime(p, time.Seconds())
}

// WriteAtTime appends len(p) bytes to the File that corresponds to time.
// If this File does not already exist, it will be created.
// It returns the number of bytes written and an Error.
func (w *Writer) WriteAtTime(p []byte, time int64) (n int, err os.Error) {
	err = w.open(time)
	if err != nil {
		return 0, err
	}
	return w.file.Write(p)
}

// WriteString is equivalent to WriteStringAtTime(s, time.Seconds())
func (w *Writer) WriteString(s string) (n int, err os.Error) {
	return w.WriteStringAtTime(s, time.Seconds())
}

// WriteStringAtTime is like WriteAtTime, but writes the contents of string s instead of an array of bytes.
func (w *Writer) WriteStringAtTime(s string, time int64) (n int, err os.Error) {
	err = w.open(time)
	if err != nil {
		return 0, err
	}
	return w.file.WriteString(s)
}

// Sync calls os.File.Sync() on the File currently open for writing.
func (w *Writer) Sync() (err os.Error) {
	if w.file != nil {
		err = w.file.Sync()
	}
	return err
}

// Close calls os.File.Close() on the File that is currently open for writing. 
func (w *Writer) Close() (err os.Error) {
	if w.file != nil {
		err = w.file.Close()
	}
	return err
}
