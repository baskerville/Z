# Description

This is a small variation around rupa's [z](https://github.com/rupa/z) and my first incursion into [Go](http://golang.org/).

# Synopsis

    Z [OPTIONS] [PATTERN ...]

# Options

- `-a ITEM` — Add the given item to the data file.
- `-d ITEM` — Delete the given item from the data file.
- `-i PATH` — Use the given file as data file.
- `-s frecency|hits|atime` — Use the given sort method.

# Environment Variables

- `Z_DATA_FILE`: path to the data file (defaults to `~/.z`).
- `Z_HISTORY_SIZE`: maximum number of items.
- `Z_AGING_CONSTANT`: value of the aging constant.

# Help

Ensure that the data file exists before running `Z`.

If you are a `z` user, you can generate a data file with:

```
awk -F '|' 'BEGIN {OFS="\0"} {print $3, int($2), $1}' ~/.z
```

You can emulate `z` with:

```
z() {
    pushd "$(Z "$@" | head -n 1)" > /dev/null
}
```

And by adding the following to your shell's `PROMPT_COMMAND`:

```
[ "$PWD" -ef "$HOME" ] || Z -a "$PWD"
```

# Frecency

The *frecency* is given by:
```
h * A / (A + t)
```
Where `h` is the number of hits, `t` the access time and `A` the aging constant.
