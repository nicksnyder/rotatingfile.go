package rotatingfile

import (
	"os"
	"time"
)

type Writer struct {
	RotatingFile
	file *os.File // The current file being written to
}

func NewWriter(secondsPerFile int, format string) *Writer {
	return &Writer{RotatingFile{secondsPerFile, format}, nil}
}

func (w *Writer) open(time int64) (err os.Error) {
	filename := w.filename(time, 0)

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

func (w *Writer) Write(p []byte) (n int, err os.Error) {
	return w.WriteAtTime(p, time.Nanoseconds())
}

func (w *Writer) WriteAtTime(p []byte, time int64) (n int, err os.Error) {
	err = w.open(time)
	if err != nil {
		return 0, err
	}
	return w.file.Write(p)
}

func (w *Writer) WriteString(s string) (n int, err os.Error) {
	return w.WriteStringAtTime(s, time.Nanoseconds())
}

func (w *Writer) WriteStringAtTime(s string, time int64) (n int, err os.Error) {
	err = w.open(time)
	if err != nil {
		return 0, err
	}
	return w.file.WriteString(s)
}

func (w *Writer) Sync() (err os.Error) {
	if w.file != nil {
		err = w.file.Sync()
	}
	return err
}

func (w *Writer) Close() (err os.Error) {
	if w.file != nil {
		err = w.file.Close()
	}
	return err
}
