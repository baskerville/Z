NAME = Z
VERSION = 0.6

PREFIX ?= /usr/local
BINPREFIX = $(PREFIX)/bin
MANPREFIX = $(PREFIX)/share/man

all: $(NAME)

$(NAME): $(NAME).go
	go build $(NAME).go

install:
	mkdir -p "$(DESTDIR)$(BINPREFIX)"
	cp -p $(NAME) "$(DESTDIR)$(BINPREFIX)"
	mkdir -p "$(DESTDIR)$(MANPREFIX)"/man1
	cp -p doc/$(NAME).1 "$(DESTDIR)$(MANPREFIX)"/man1

uninstall:
	rm -f "$(DESTDIR)$(BINPREFIX)"/$(NAME)
	rm -f "$(DESTDIR)$(MANPREFIX)"/man1/$(NAME).1

doc:
	a2x -v -d manpage -f manpage -a revnumber=$(VERSION) doc/$(NAME).1.txt

clean:
	rm -f $(NAME)

.PHONY: all $(NAME) install uninstall doc clean 
