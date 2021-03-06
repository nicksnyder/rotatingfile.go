// You can run the example with this command (assuming 6g compiler): 
// 6g example.go && 6l example.6 && time ./6.out && rm example.*.log;
package main

import (
	"fmt"
	"github.com/nicksnyder/rotatingfile.go"
	"time"
)

const (
	secondsPerFile       = 2
	format               = "example.%v.log"
	nanosecondsPerSecond = 1000000000
)

func main() {
	beginTime := time.Seconds()
	w := rotatingfile.NewWriter(secondsPerFile, format)
	w.WriteString("hello\n")
	time.Sleep(secondsPerFile * nanosecondsPerSecond)
	w.WriteString("world\n")
	endTime := time.Seconds()

	r := rotatingfile.NewReader(secondsPerFile, format, beginTime, endTime)
	buf := make([]byte, 12)
	_, err := r.Read(buf)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", buf)
}
