include ../Make.inc

TARG=metricsd
GOFILES=\
	main.go\
	cli.go
include $(GOROOT)/src/Make.cmd

start: all
	./metricsd -debug=0 -slice=1 -write=5

godoc:
	open http://localhost:6060/
	godoc -path=. -http=:6060
