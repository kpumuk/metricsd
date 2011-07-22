include ../Make.inc

TARG=gorrdpd
GOFILES=main.go
include $(GOROOT)/src/Make.cmd

start: all
	./gorrdpd -debug=0 -slice=1 -write=5

godoc:
	open http://localhost:6060/
	godoc -path=. -http=:6060
# 
# clean:
# 	rm -rf *.o *.a *.[$(OS)] [$(OS)].out $(CLEANFILES)
# 	make -C stdlib clean
