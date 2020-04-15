# Filme [![Build Status](https://travis-ci.org/florinutz/filme.svg?branch=master)](https://travis-ci.org/florinutz/filme) [![codecov](https://codecov.io/gh/florinutz/filme/branch/master/graph/badge.svg)](https://codecov.io/gh/florinutz/filme)

### A movie torrenting utility.

Filme is a little torrenting helper I started writing for myself some time ago.
The one solid feature it has right now is the torrent searching, which
basically crawls 1337x.to from the cli.

It can search for torrents (`filme help search`) and display details about any of the search results (`filme help detail`).

## Install

### snap

[![Get it from the Snap Store](https://snapcraft.io/static/images/badges/en/snap-store-white.svg)](https://snapcraft.io/filme)

```bash
snap install filme
```

### build from source

Run `make` as an entry point:

```bash
make binary
sudo mv build/filme-linux-amd64 /usr/local/bin/filme
```

## Usage

#### Search
List all items matching "1080p". Crawl first 20 pages:

```bash
filme search 1080p --page-min 1 --page-max 20
```

Please crawl 1337x politely by using the `--reqs-*` crawler settings. They set delays and parallelism.
You can also set the user agent, but leaving it empty will use valid random ones from different browsers.

#### Display details

Display one of the items:

```bash
filme detail https://1337x.to/torrent/4287659/Honey-Boy-2019-1080p-WEBRip-5-1-YTS-YIFY/
```
will get you
```
Gemini Man (2019) [WEBRip] [1080p] [YTS] [YIFY]

id: 4129693
magnet: magnet:?xt=urn:btih:671E6D130005810236C19F2B706AB5552CA1472A&dn=Gemini+Man+%282019%29+%5BW...

seeders: 14594
leechers: 7396
```

You can also get its imdb link if you have an [OMDB api key](https://www.omdbapi.com/apikey.aspx) 
into the `OMDB_API_KEY` env var:
```bash
OMDB_API_KEY=84824e9a filme detail https://1337x.to/torrent/4287659/Honey-Boy-2019-1080p-WEBRip-5-1-YTS-YIFY/
```

```
...
Imdb info:
https://www.imdb.com/title/tt1025100/ (2019) Gemini Man
``` 

## Caching
Everything is cached in `~/.cache/filme` and there is currently no way to avoid this, so please delete this folder 
whenever you want fresh results.

```bash
rm -rf ~/.cache/filme
```

![Keybase PGP](https://img.shields.io/keybase/pgp/fl0?style=social)
![Keybase XLM](https://img.shields.io/keybase/xlm/fl0?style=social)