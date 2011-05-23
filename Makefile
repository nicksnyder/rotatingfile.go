include $(GOROOT)/src/Make.inc

GOFMT=gofmt
TARG=rotatingfile
GOFILES=\
	file.go\
	reader.go\
	writer.go\

include $(GOROOT)/src/Make.pkg

format:
	${GOFMT} -w *.go
