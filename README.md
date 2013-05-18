# Description

This is a small variation around [rupa](https://github.com/rupa)'s [z](https://github.com/rupa/z) and my first incursion into [Go](http://golang.org/).

# Synopsis

    Z [OPTIONS] [PATTERN ...]

# Options

- `-a ITEM` — Add the given item to the data file.
- `-d ITEM` — Delete the given item from the data file.
- `-i PATH` — Use the given file as data file.
- `-s frecency|hits|atime` — Use the given sort method.

# Environment Variables

- `Z_DATA_FILE`: path to the data file.
- `Z_HISTORY_SIZE`: maximum number of items.
- `Z_AGING_CONSTANT`: value of the aging constant.
