include $(GOROOT)/src/Make.inc

TARG=rotatingfile
GOFILES=\
	file.go\
	reader.go\
	writer.go\

include $(GOROOT)/src/Make.pkg
