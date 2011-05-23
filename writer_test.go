package rotatingfile

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"testing"
)

var (
	line = []byte("test\n")
	sum  = []byte{0xb4, 0x8b, 0x35, 0x35, 0xcd, 0x5e, 0x1e, 0xb7, 0xdc, 0xd4, 0x4c, 0x79, 0x5e, 0x76, 0x4e, 0x4f, 0x26, 0x1e, 0xde, 0x31}
)

func TestWrite(t *testing.T) {
	dir := getTempDir(t)
	defer os.RemoveAll(dir)

	format := getFormat(dir)

	// Write log files
	w := NewWriter(secondsPerFile, format)
	defer w.Close()
	for i := 0; i < seconds; i++ {
		_, err := w.WriteAtTime(line, int64(i))
		if err != nil {
			t.Fatal(err)
		}
	}

	// Verify log files
	hash := sha1.New()
	for i := 0; i < seconds; i += secondsPerFile {
		filename := fmt.Sprintf(format, i)
		f, err := os.Open(filename)
		if err != nil {
			t.Fatal(err)
		}
		hash.Reset()
		io.Copy(hash, f)
		if !bytes.Equal(hash.Sum(), sum) {
			t.Fatalf("Hash of %s was %x, expected %x", filename, hash.Sum(), sum)
		}
	}
}
