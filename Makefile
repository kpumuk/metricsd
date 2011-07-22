ifeq ($(DESTINATION),)
DESTINATION:=/usr/local/gorrdpd
endif

all: build

build:
	GOPATH=$(CURDIR) goinstall -clean gorrdpd
	GOPATH=$(CURDIR) goinstall -clean benchmark

install: build rrdtool
	mkdir -p $(DESTINATION)/data
	if test -e $(DESTINATION)/gorrdpd.old; \
	then rm -f $(DESTINATION)/gorrdpd.old; \
	fi
	if test -e $(DESTINATION)/gorrdpd; \
	then mv $(DESTINATION)/gorrdpd $(DESTINATION)/gorrdpd.old; \
	fi
	cp -r bin/gorrdpd bin/gorrdpd.sh templates $(DESTINATION)
	if test ! -e $(DESTINATION)/gorrdpd.conf; \
	then cp gorrdpd.conf.sample $(DESTINATION)/gorrdpd.conf; \
	fi

rrdtool:
	if test ! -e /usr/bin/rrdtool; \
	then echo "Please install rrdtool to /usr/bin/rrdtool"; exit; \
	fi

clean:
	if test -e $(GORRD_DIR); \
	then make -C $(GORRD_DIR) clean; \
	fi
	make -C src clean
