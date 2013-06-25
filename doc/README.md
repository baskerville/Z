% Z User's Guide
% Bastien Dejean
% June 25, 2013

# Synopsis

**Z** [*OPTIONS*] [*PATTERN* ...]

# Description

This is a small variation around rupa's [z](https://github.com/rupa/z).

# Options

`-a` *ITEM*
:    Add the given item to the data file.
`-d` *ITEM*
:    Delete the given item from the data file.
`-i` *PATH*
:    Use the given file as data file.
`-s` *frecency*|*hits*|*atime*
:    Use the given sort method (defaults to *frecency*).

# Environment Variables

*Z\_DATA\_FILE*
:    path to the data file (defaults to *~/.z*).
*Z\_HISTORY\_SIZE*
:    maximum number of items.
*Z\_AGING\_CONSTANT*
:    value of the aging constant.

# Help

Ensure that the data file exists before running **Z**.

Either create an empty one:

```
touch ~/.z
```

Or, if you are a **z** user, import your data file with:

```
mv ~/.z ~/.z.bak
awk -F '|' 'BEGIN {OFS="\0"} {print $3, int($2), $1}' ~/.z.bak > ~/.z
```

You can emulate **z** with:

```
z() {
    local dir=$(Z "$@" | head -n 1)
    pushd "$dir" > /dev/null 2>&1 || Z -d "$dir"
}
```

If your shell is **Bash**, add the following to *~/.bashrc*:

```
export PROMPT_COMMAND='[ "$PWD" -ef "$HOME" ] || Z -a "$PWD"'
```

Else, if your shell is **Zsh**, add the following to *~/.zshrc*:

```
chpwd() {
    [ "$PWD" -ef "$HOME" ] || Z -a "$PWD"
}
```

# Frecency

The *frecency* is given by:

```
h * A / (A + t)
```

Where *h* is the number of hits, *t* the access time and *A* the aging constant.
