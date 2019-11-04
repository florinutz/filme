# Filme [![Build Status](https://travis-ci.org/florinutz/filme.svg?branch=master)](https://travis-ci.org/florinutz/filme) [![codecov](https://codecov.io/gh/florinutz/filme/branch/master/graph/badge.svg)](https://codecov.io/gh/florinutz/filme)

### A movie torrenting utility.

Filme exposes a set of crawlers meant to be used to parse data related to movie torrents. The initial idea was
to expose a unified tool that lists latest movie torrents next to their imdb rating and rottentomatoes reviews,
but it was meant to fail, since such a solution will always be illegal.

I will leave these here for anyone interested in continuing my work.

Run `make` as an entry point:

```bash
make binary # creates a ./bin symlink
make test
...
./bin crawl google independence day
./bin crawl 1337x_detail https://1337x.to/torrent/3746692/Spaced-Invaders-1990-BluRay-720p-YTS-YIFY/
./bin crawl 1337x_list
```

Please contribute! (add more crawlers?)