NAME = Z

PREFIX ?= /usr/local
BINPREFIX = $(PREFIX)/bin
MANPREFIX = $(PREFIX)/share/man

all: $(NAME)

$(NAME): $(NAME).go
	go build $(NAME).go

install:
	mkdir -p "$(DESTDIR)$(BINPREFIX)"
	cp -p $(NAME) "$(DESTDIR)$(BINPREFIX)"
	cp -p $(NAME).1 "$(DESTDIR)$(MANPREFIX)"/man1

uninstall:
	rm -f "$(DESTDIR)$(BINPREFIX)"/$(NAME)
	rm -f "$(DESTDIR)$(MANPREFIX)"/man1/$(NAME).1
doc:
	pandoc -t json doc/README.md | runhaskell doc/man_filter.hs | pandoc --no-wrap -f json -t man --template doc/man.template -V name=$(NAME) -o $(NAME).1
	pandoc --no-wrap -f markdown -t asciidoc doc/README.md -o README.asciidoc
	patch -p 1 -i doc/quirks.patch
clean:
	rm -f $(NAME)

.PHONY: all $(NAME) install uninstall doc clean 
