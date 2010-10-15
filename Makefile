ifeq ($(DESTINATION),)
DESTINATION:=/usr/local/gorrdpd
endif
GORRD_DIR=gorrd.git

all: build

gorrd:
	if test ! -e $(GORRD_DIR); \
	then git clone -q git@github.com:kpumuk/gorrd.git $(GORRD_DIR); \
	else cd $(GORRD_DIR) && git pull -q; \
	fi
	make -C $(GORRD_DIR) install

web.go:
	goinstall github.com/hoisie/web.go

mustache.go:
	goinstall github.com/hoisie/mustache.go

rrdtool:
	if test ! -e /usr/bin/rrdtool; \
	then echo "Please install rrdtool to /usr/bin/rrdtool"; exit; \
	fi

build: gorrd web.go mustache.go
	make -C src gorrdpd

install: build rrdtool
	mkdir -p $(DESTINATION)/data
	if test -e $(DESTINATION)/gorrdpd.old; \
	then rm -f $(DESTINATION)/gorrdpd.old; \
	fi
	if test -e $(DESTINATION)/gorrdpd; \
	then mv $(DESTINATION)/gorrdpd $(DESTINATION)/gorrdpd.old; \
	fi
	cp -r src/gorrdpd script/gorrdpd.sh src/templates $(DESTINATION)
	if test ! -e $(DESTINATION)/gorrdpd.conf; \
	then cp src/gorrdpd.conf.sample $(DESTINATION)/gorrdpd.conf; \
	fi

clean:
	if test -e $(GORRD_DIR); \
	then make -C $(GORRD_DIR) clean; \
	fi
	make -C src clean
