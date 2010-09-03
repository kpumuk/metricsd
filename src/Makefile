include $(GOROOT)/src/Make.inc

TARG=gorrdpd
GOFILES=\
	server.go
PREREQ=\
	stdlib.make \
	logger.$O \
	config.$O \
	types.$O \
	parser.$O \
	writers.$O
WRITERS=$(wildcard writers/*.go)

include $(GOROOT)/src/Make.cmd

%.$O: %.go
	${GC} $<

writers.$O: writers.go $(WRITERS)
	${GC} writers.go $(WRITERS)

start: all
	./gorrdpd -debug=0 -slice=1 -write=5

godoc:
	open http://localhost:6060/
	godoc -path=. -http=:6060
