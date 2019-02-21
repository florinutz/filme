package coll33tx

import (
	"net"
	"net/http"
	"time"

	"github.com/florinutz/filme/pkg/collector"

	"github.com/gocolly/colly"
	collyExtensions "github.com/gocolly/colly/extensions"
	log "github.com/sirupsen/logrus"
)

// DomainConfig is a colly extension that configures both 1337x collectors
func DomainConfig(c *colly.Collector, logEntry *log.Entry) {
	for _, f := range []func(collector *colly.Collector){
		colly.MaxDepth(1),
		colly.Async(true),
		colly.CacheDir(".cache"),
		colly.UserAgent("filme finder"),
		colly.AllowedDomains("1337x.to"),
		colly.Debugger(&collector.LogrusDebugger{Logger: logEntry.Logger}),
	} {
		f(c)
	}

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
}
