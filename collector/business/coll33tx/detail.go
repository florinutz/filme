package coll33tx

import (
	"github.com/florinutz/filme/collector/coll33tx"
	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

// NewDetailCollector creates a new 1337x.to details page collector tweaked for business
func NewDetailCollector(log *logrus.Entry, options ...func(collector *colly.Collector)) *coll33tx.DetailsCollector {
	return coll33tx.NewDetailsPageCollector(torrentFound(log), options...)
}

func torrentFound(log *logrus.Entry) func(torrent coll33tx.L33tTorrent) {
	return func(torrent coll33tx.L33tTorrent) {
		log.WithField("torrent", torrent).Debug("torrent found on detail page")
	}
}
