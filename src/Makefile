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
	writers.$O \
	web.$O
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

clean:
	rm -rf *.o *.a *.[$(OS)] [$(OS)].out $(CLEANFILES)
	make -C stdlib clean

test:
	gotest

testpackage-clean:
	rm -f _test/$(TARG).a _gotest_.$O

testpackage: _test/parser.a

_gotest_.$O: parser.go parser_test.go types.$O
	$(GC) -o $@ parser.go parser_test.go

_test/parser.a: _gotest_.$O $(OFILES)
	@mkdir -p _test/$(dir)
	rm -f _test/parser.a
	gopack grc $@ _gotest_.$O $(OFILES)

importpath:
	@echo parser
