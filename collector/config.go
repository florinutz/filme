package collector

import (
	"time"

	"net"
	"net/http"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

const (
	CtxKeyMagnet = "magnet"
	CtxKeyTitle  = "title"
)

// Configure adds common configs for all filme Collectors
func Configure(c *colly.Collector) {
	for _, option := range []func(collector *colly.Collector){
		colly.MaxDepth(2),
		colly.Async(true),
		colly.CacheDir("/tmp/cache"),
	} {
		option(c)
	}

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

	// Rotate socks5 proxies
	//rp, err := proxy.RoundRobinProxySwitcher("sockss5://wOBzsRUmerF:A7691RHzprQ@ams.socks.ipvanish.com")
	//if err != nil {
	//	log.WithError(err).Fatal("Couldn't use the socks5 proxy")
	//}
	//c.SetProxyFunc(rp)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})

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
		}).Warn("Crawling error!")
	})
}
