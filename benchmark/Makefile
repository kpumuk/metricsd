include $(GOROOT)/src/Make.$(GOARCH)

TARG=udp_generator
GOFILES=\
	udp_generator.go \

include $(GOROOT)/src/Make.pkg

all: udp_generator

%.$O: %.go
	${GC} $<

udp_generator: udp_generator.$O
	${LD} -o $@ udp_generator.$O

start: udp_generator
	./udp_generator
