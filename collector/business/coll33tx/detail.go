package coll33tx

import (
	"encoding/json"

	"fmt"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"gitlab.com/phlo/filme/collector/coll33tx"
)

// NewDetailCollector creates a new 1337x.to details page collector tweaked for business
func NewDetailCollector(
	log *logrus.Entry,
	options ...func(collector *colly.Collector),
) *coll33tx.DetailsCollector {
	collector := coll33tx.NewDetailsPageCollector(torrentFound(log), options...)
	return collector
}

func torrentFound(log *logrus.Entry) func(torrent coll33tx.L33tTorrent) {
	return func(torrent coll33tx.L33tTorrent) {
		log := log.WithField("torrent", torrent)
		log.Debug("torrent found on detail page")

		j, err := json.Marshal(torrent)
		if err != nil {
			log.Error("could not encode torrent to json")
		}

		fmt.Println(string(j))
	}
}
