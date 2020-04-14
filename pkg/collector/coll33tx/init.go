package coll33tx

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/florinutz/filme/pkg/collector"

	"github.com/gocolly/colly"
	collyExtensions "github.com/gocolly/colly/extensions"
	log "github.com/sirupsen/logrus"
)

// DomainConfig is a colly extension that configures both 1337x collectors
func DomainConfig(c *colly.Collector, delay, randomDelay, parallelism int, userAgent string, logEntry log.Entry) {
	cacheDir := getCacheDir()

	logEntry.WithField("cache_dir", cacheDir).Debug("using cache dir")

	for _, f := range []func(collector *colly.Collector){
		colly.MaxDepth(1),
		colly.Async(true),
		colly.CacheDir(cacheDir),
		colly.AllowedDomains("1337x.to"),
		colly.Debugger(&collector.LogrusDebugger{Logger: logEntry.Logger}),
	} {
		f(c)
	}

	if userAgent == "" {
		collyExtensions.RandomUserAgent(c)
	} else {
		colly.UserAgent(userAgent)(c)
	}

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
		Delay:       time.Duration(delay) * time.Second,
		RandomDelay: time.Duration(randomDelay) * time.Second,
		Parallelism: parallelism,
		DomainGlob:  "1337x.to",
	}); err != nil {
		logEntry.WithError(err).Fatal("failed while setting collector limit")
	}
}

func getCacheDir() string {
	p1 := fmt.Sprintf("%s/.cache", homeDir())
	if _, err := os.Stat(p1); os.IsNotExist(err) {
		_ = os.Mkdir(p1, 0750)
	}
	p2 := fmt.Sprintf("%s/filme", p1)
	if _, err := os.Stat(p2); os.IsNotExist(err) {
		_ = os.Mkdir(p2, 0750)
	}

	return p2
}

func homeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
