package coll33tx

import (
	"net"
	"net/http"
	"time"

	debug "github.com/florinutz/filme/collector"

	"github.com/gocolly/colly"
	collyExtensions "github.com/gocolly/colly/extensions"
	log "github.com/sirupsen/logrus"
)

func initCollector(logEntry *log.Entry, options ...func(collector *colly.Collector)) *colly.Collector {
	c := colly.NewCollector(append(options,
		colly.MaxDepth(1),
		colly.Async(true),
		colly.CacheDir(".cache"),
		colly.UserAgent("filme finder"),
		colly.AllowedDomains("1337x.to"),
		colly.Debugger(&debug.LogrusDebugger{Logger: logEntry.Logger}),
	)...)

	collyExtensions.RandomUserAgent(c)
	collyExtensions.Referrer(c)

	c.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   7 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})

	if err := c.Limit(&colly.LimitRule{
		Delay:       3 * time.Second,
		RandomDelay: 3 * time.Second,
		DomainGlob:  "1337x.to",
		Parallelism: 4,
	}); err != nil {
		logEntry.WithError(err).Fatal("failed while setting collector limit")
	}

	return c
}
