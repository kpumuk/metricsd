# figure out what GOROOT is supposed to be
GOROOT ?= $(shell printf 't:;@echo $$(GOROOT)\n' | gomake -f -)
include $(GOROOT)/src/Make.inc

# format:
# 	gofmt -spaces=true -tabindent=false -tabwidth=4 -s -w $(wildcard *.go)
