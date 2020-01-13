# Filme [![Build Status](https://travis-ci.org/florinutz/filme.svg?branch=master)](https://travis-ci.org/florinutz/filme) [![codecov](https://codecov.io/gh/florinutz/filme/branch/master/graph/badge.svg)](https://codecov.io/gh/florinutz/filme)

### A movie torrenting utility.

Filme is a little torrenting helper I started writing for myself some time ago.
The one solid feature it has right now is the torrent searching, which
basically crawls 1337x.to from the cli.

I've planned a lot more features.

## Install

### snap
```bash
# I'm still working on this, so please build from source
snap install filme
```

### build from source

Run `make` as an entry point:

```bash
make binary
./bin help search
```

## Usage

```bash
./bin help search 
Search torrents

Usage:
  filme search <what> [flags]

Flags:
  -c, --category category   one of: all, movies, tv, documentaries, anime, xxx (default all)
  -d, --crawl-details       follows every link in the list and fetches detail pages data
  -e, --encoding encoding   one of: dvd, h264-x264, hd, uhd, hevc-x265, mp4, svcd-vcd, divx-xvid
  -h, --help                help for search
      --leechers-max uint   ignores items with more leechers
      --leechers-min uint   ignores items with less leechers
      --page-max uint       stop at this page
      --page-min uint       start at this page
      --seeders-max uint    ignores items with more seeders
      --seeders-min uint    ignores items with less seeders
      --size uint           ignore items bigger than this
  -s, --sort sort-pair      one of: time-asc, time-desc, size-asc, size-desc, seeders-asc, seeders-desc, leechers-asc, leechers-desc (default seeders-desc)
  -t, --total uint          specifies the maximum desired number of items to display.
                            Defaults to one page's worth of items. (default 20)

Global Flags:
      --debug-level level     one of: panic, fatal, error, warning, info, debug, trace (default panic)
      --debug-report-caller   show debug callers
```
## Examples

### search 

Please use and contribute!
