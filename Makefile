PREFIX ?= /usr/local
BINPREFIX = $(PREFIX)/bin

all: Z

Z: Z.go
	go build Z.go

install:
	mkdir -p "$(DESTDIR)$(BINPREFIX)"
	cp -p Z "$(DESTDIR)$(BINPREFIX)"
clean:
	rm -f Z

.PHONY: all Z install clean 
