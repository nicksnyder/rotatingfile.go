package rotatingfile

import (
	"io/ioutil"
	"testing"
	"path/filepath"
)

var (
	format = "%d.log"
	seconds        = 100
	secondsPerFile = 10
)

type timeTest struct {
	time, index, newTime int64
}

type filenameTest struct {
	time, index int64
	filename string
}

var timeTests = []timeTest{
	timeTest{0, 0, 0},
	timeTest{0, 1, 10},
	timeTest{0, 2, 20},
	timeTest{5, 0, 0},
	timeTest{5, 1, 10},
	timeTest{5, 2, 20},
	timeTest{10, 0, 10},
	timeTest{10, 1, 20},
	timeTest{10, 2, 30},
	timeTest{50, 0, 50},
	timeTest{50, 1, 60},
	timeTest{50, 2, 70},
}

var filenameTests = []filenameTest{
	filenameTest{0, 0, "0.log"},
	filenameTest{0, 1, "10.log"},
	filenameTest{0, 2, "20.log"},
	filenameTest{5, 0, "0.log"},
	filenameTest{5, 1, "10.log"},
	filenameTest{5, 2, "20.log"},
}

func getTempDir(t *testing.T) (dir string) {
	dir, err := ioutil.TempDir("", "rotatingfile")
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func getFormat(dir string) string {
	return filepath.Join(dir, format)
}

func TestTime(t *testing.T) {
	rf := &RotatingFile{secondsPerFile, format}
	for _, tt := range timeTests {
		newTime := rf.time(tt.time, tt.index)
		if newTime != tt.newTime {
			t.Errorf("time(%d, %d) = %d, expected %d", tt.time, tt.index, newTime, tt.newTime)
		}
	}
}

func TestFilename(t *testing.T) {
	rf := &RotatingFile{secondsPerFile, format}
	for _, ft := range filenameTests {
		filename := rf.filename(ft.time, ft.index)
		if filename != ft.filename {
			t.Errorf("filename(%d, %d) = %s, expected %s", ft.time, ft.index, filename, ft.filename)
		}
	}
}
