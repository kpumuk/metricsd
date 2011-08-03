# figure out what GOROOT is supposed to be
GOROOT ?= $(shell printf 't:;@echo $$(GOROOT)\n' | gomake -f -)
include $(GOROOT)/src/Make.inc

GCIMPORTS=-I "$(GOPATH)/pkg/$(GOOS)_$(GOARCH)"
export LDIMPORTS=-L _test -L $(GOPATH)/pkg/$(GOOS)_$(GOARCH)/
export GL=$(GOPATH)/src/gl.sh
