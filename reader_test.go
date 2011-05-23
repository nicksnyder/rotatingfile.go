package rotatingfile

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

var (
	begin int64 = 0
	end int64 = 100
	r *Reader = nil
)

type seekToTimeTest struct {
	time, offset, newTime, newOffset int64
	err os.Error
}

type readTest struct {
	time, offset int64
	buf []byte
	err os.Error
}

var seekToTimeTests = []seekToTimeTest{
	seekToTimeTest{0, 0, 10, 0, nil},
	seekToTimeTest{0, 5, 10, 5, nil},
	seekToTimeTest{0, 10, 20, 0, nil},
	seekToTimeTest{0, 35, 50, 5, nil},
	seekToTimeTest{0, 36, end + int64(secondsPerFile), 0, os.EOF},
	seekToTimeTest{0, 37, end + int64(secondsPerFile), 1, os.EOF},
	seekToTimeTest{5, 0, 10, 0, nil},
	seekToTimeTest{45, 15, 50, 5, nil},
	seekToTimeTest{50, 5, 50, 5, nil},
	seekToTimeTest{50, 6, end + int64(secondsPerFile), 0, os.EOF},
	seekToTimeTest{50, 7, end + int64(secondsPerFile), 1, os.EOF},
	seekToTimeTest{50, -1, 40, 9, nil},
	seekToTimeTest{99, 0, end + int64(secondsPerFile), 0, os.EOF},
	seekToTimeTest{99, -1, 50, 5, nil},
	seekToTimeTest{99, -6, 50, 0, nil},
	seekToTimeTest{99, -7, 40, 9, nil},
	seekToTimeTest{99, -35, 10, 1, nil},
	seekToTimeTest{99, -36, 10, 0, nil},
	seekToTimeTest{99, -37, begin - int64(secondsPerFile), -1, os.EOF},
}

var readTests = []readTest{
	readTest{0, 0, []byte("0123456789abcdefghijklmnopqrstuvwxyz"), nil},
	readTest{5, 5, []byte("56789abcde"), nil},
	readTest{19, 5, []byte("56789abcde"), nil},
	readTest{20, 15, []byte("pqrstuvwxyz"), nil},
	readTest{50, 0, []byte("uvwxyz\x00"), os.EOF},
}

func init() {
	filename := filepath.Join("testdata", format)
	r = NewReader(secondsPerFile, filename, begin, end)
}

func TestSeekToTime(t *testing.T) {
	for _, st := range seekToTimeTests {
		newTime, newOffset, err := r.SeekToTime(st.time, st.offset)
		if newTime != st.newTime || newOffset != st.newOffset || err != st.err {
			t.Errorf("SeekToTime(%d, %d) got time=%d offset=%d err=%v, expected time=%d offset=%d err=%v", st.time, st.offset, newTime, newOffset, err, st.newTime, st.newOffset, st.err)
		}
	}
}

func TestRead(t *testing.T) {
	for _, rt := range readTests {
		_, _, err := r.SeekToTime(rt.time, rt.offset)
		if err != nil {
			t.Errorf("SeekToTime(%d, %d) got err=%v, expected err=nil", err)
		}
		buf := make([]byte, len(rt.buf))
		_, err = r.Read(buf)
		if bytes.Compare(buf, rt.buf) != 0 || err != rt.err {
			t.Errorf("Read() returned buf=%s err=%v, expected buf=%s err=%v", buf, err, rt.buf, rt.err)
		}
	}
}

func TestReadAt(t *testing.T) {
	for _, rt := range readTests {
		buf := make([]byte, len(rt.buf))
		_, err := r.ReadAt(buf, rt.time, rt.offset)
		if bytes.Compare(buf, rt.buf) != 0 || err != rt.err {
			t.Errorf("ReadAt() returned buf=%s err=%v, expected buf=%s err=%v", buf, err, rt.buf, rt.err)
		}
	}
}
