include $(GOROOT)/src/Make.inc

TARG=gorrdpd
GOFILES=\
	server.go
PREREQ=\
	logger.$O \
	config.$O \
	types.$O \
	writers.$O \

include $(GOROOT)/src/Make.cmd

%.$O: %.go
	${GC} $<

start: all
	./gorrdpd -debug=0 -slice=1 -write=5

godoc:
	open http://localhost:6060/
	godoc -path=. -http=:6060
