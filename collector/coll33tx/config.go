package coll33tx

import (
	"net"
	"net/http"
	"time"

	"github.com/gocolly/colly"
	collyExtensions "github.com/gocolly/colly/extensions"
	log "github.com/sirupsen/logrus"
)

func initCollector(options ...func(collector *colly.Collector)) *colly.Collector {
	c := colly.NewCollector(append(options,
		colly.MaxDepth(1),
		colly.Async(true),
		colly.CacheDir(".cache"),
	)...)

	c.AllowedDomains = []string{"1337x.to"}
	c.UserAgent = "filme finder"

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
