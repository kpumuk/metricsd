ifeq ($(DESTINATION),)
DESTINATION:=/usr/local/metricsd
endif

all: build

build:
	git submodule update
	GOPATH=$(CURDIR) goinstall -clean metricsd
	GOPATH=$(CURDIR) goinstall -clean benchmark

install: build rrdtool
	mkdir -p $(DESTINATION)/data
	if test -e $(DESTINATION)/metricsd.old; \
	then rm -f $(DESTINATION)/metricsd.old; \
	fi
	if test -e $(DESTINATION)/metricsd; \
	then mv $(DESTINATION)/metricsd $(DESTINATION)/metricsd.old; \
	fi
	cp -r bin/metricsd bin/metricsd.sh templates public $(DESTINATION)
	if test ! -e $(DESTINATION)/metricsd.conf; \
	then cp metricsd.conf.sample $(DESTINATION)/metricsd.conf; \
	fi

rrdtool:
	if test ! -e /usr/bin/rrdtool; \
	then echo "Please install rrdtool to /usr/bin/rrdtool"; exit; \
	fi

clean:
	make -C src clean
