include $(GOROOT)/src/Make.$(GOARCH)

TARG=gorrdpd
GOFILES=\
	server.go \

server.$O: config.$O logger.$O types.$O writers.$O
config.$O: logger.$O

include $(GOROOT)/src/Make.pkg

all: gorrdpd

%.$O: %.go
	${GC} $<

gorrdpd: server.$O
	${LD} -o $@ server.$O

start: gorrdpd
	./gorrdpd -debug=0 -slice=1 -write=5 -data ./data
