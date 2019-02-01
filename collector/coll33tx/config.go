package coll33tx

import (
	"net"
	"net/http"
	"time"

	"github.com/gocolly/colly"
	collyExtensions "github.com/gocolly/colly/extensions"
	log "github.com/sirupsen/logrus"
)

func getCollyCollector(options ...func(collector *colly.Collector)) *colly.Collector {
	c := colly.NewCollector(options...)

	c.AllowedDomains = []string{"1337x.to"}
	c.UserAgent = "filme finder"

	collyExtensions.RandomUserAgent(c)
	collyExtensions.Referrer(c)

	colly.MaxDepth(1)(c)
	colly.Async(true)(c)
	colly.CacheDir(".cache")(c)

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
		DomainGlob:  "*",
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	}); err != nil {
		log.WithError(err).Fatal("failed while setting collector limit")
	}

	c.OnRequest(func(r *colly.Request) {
		log.WithFields(log.Fields{"url": r.URL}).Debug("Requesting a url")
	})

	c.OnResponse(func(r *colly.Response) {
		log.WithFields(log.Fields{
			"url":     r.Request.URL,
			"content": string(r.Body),
		}).Debug("Got a response")
	})

	c.OnError(func(r *colly.Response, err error) {
		log.WithError(err).WithFields(log.Fields{
			"url":      r.Request.URL,
			"response": *r,
		}).Warn("crawling error")
	})

	return c
}
