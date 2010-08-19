ifeq ($(DESTINATION),)
DESTINATION:=/usr/local/gorrdpd
endif

all: build

build:
	make -C src gorrdpd

install: build
	mkdir -p $(DESTINATION)/data
	cp src/gorrdpd script/gorrdpd.sh $(DESTINATION)
	if test ! -e $(DESTINATION)/gorrdpd.conf; \
	then cp src/gorrdpd.conf.sample $(DESTINATION)/gorrdpd.conf; \
	fi

clean:
	make -C src clean
