package collector

import "net/url"

type Torrents []*Torrent

type Torrent struct {
	Title   string
	Magnet  string
	FoundOn *url.URL
}
