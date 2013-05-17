# Description

This is a small variation around [rupa](https://github.com/rupa)'s [z](https://github.com/rupa/z) and my first incursion into [Go](http://golang.org/).

# Synopsis

    z [OPTIONS] [PATTERN ...]

# Options

- `-a ITEM` — Add the given item to the data file.
- `-d ITEM` — Delete the given item from the data file.
- `-i PATH` — Use the given file as data file.
- `-g VALUE` — Set the value of the aging constant.

# Environment Variables

- `Z_DATA_FILE`: path to the data file.
- `Z_HISTORY_SIZE`: maximum number of items.
