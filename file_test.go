package rotatingfile

import (
	"io/ioutil"
	"testing"
	"path/filepath"
)

var (
	format         = "%d.log"
	seconds        = 100
	secondsPerFile = 10
)

type timeTest struct {
	time, newTime int64
}

type filenameTest struct {
	time     int64
	filename string
}

var timeTests = []timeTest{
	{0, 0},
	{1, 0},
	{9, 0},
	{10, 10},
	{11, 10},
	{99, 90},
}

var filenameTests = []filenameTest{
	{0, "0.log"},
	{1, "0.log"},
	{9, "0.log"},
	{10, "10.log"},
	{11, "10.log"},
	{99, "90.log"},
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
		newTime := rf.time(tt.time)
		if newTime != tt.newTime {
			t.Errorf("time(%d) = %d, expected %d", tt.time, newTime, tt.newTime)
		}
	}
}

func TestFilename(t *testing.T) {
	rf := &RotatingFile{secondsPerFile, format}
	for _, ft := range filenameTests {
		filename := rf.Filename(ft.time)
		if filename != ft.filename {
			t.Errorf("filename(%d) = %s, expected %s", ft.time, filename, ft.filename)
		}
	}
}
