package google

import (
	"net"
	"net/http"
	"time"

	"github.com/florinutz/filme/pkg/collector"

	"github.com/gocolly/colly"
	collyExtensions "github.com/gocolly/colly/extensions"
	log "github.com/sirupsen/logrus"
)

const DomainGoogle = "www.google.com"

func DomainConfig(c *colly.Collector, logEntry *log.Entry) {

	for _, f := range []func(collector *colly.Collector){
		colly.MaxDepth(1),
		colly.Async(true),
		colly.CacheDir(".cache"),
		colly.UserAgent("filme finder"),
		colly.AllowedDomains(DomainGoogle),
		colly.Debugger(&collector.LogrusDebugger{Logger: logEntry.Logger}),
	} {
		f(c)
	}

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept-Language", "en-US,ro;q=0.8,es;q=0.5,fr;q=0.3")
	})

	collyExtensions.RandomUserAgent(c)
	collyExtensions.Referer(c)

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
		DomainGlob:  DomainGoogle,
		Parallelism: 4,
	}); err != nil {
		logEntry.WithError(err).Fatal("failed while setting collector limit")
	}
}
