include $(GOROOT)/src/Make.$(GOARCH)

TARG=gorrdpd
GOFILES=\
	server.go \

server.$O: types.$O writers.$O config.$O

include $(GOROOT)/src/Make.pkg

all: gorrdpd

%.$O: %.go
	${GC} $<

gorrdpd: server.$O
	${LD} -o $@ server.$O

start: gorrdpd
	./gorrdpd -debug=5 -slice=1 -write=5 -data ./data
