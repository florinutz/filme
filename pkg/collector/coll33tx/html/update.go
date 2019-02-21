package main

import (
	"log"
	"net/url"
	"time"

	"github.com/florinutz/filme/pkg/collector/url_mock"
)

const dataFile = "data.json"

func main() {
	url1, _ := url.Parse("https://1337x.to/torrent/3570061/House-Party-1990-WEBRip-1080p-YTS-YIFY/")
	url2, _ := url.Parse("https://1337x.to/search/romania/3/")
	url3, _ := url.Parse("https://1337x.to/popular-movies")
	urls := []*url.URL{url1, url2, url3}

	timeout := 10 * time.Second

	err := url_mock.Fetch(urls, timeout)
	if err != nil {
		log.Fatal(err)
	}

	err = url_mock.Persist(dataFile)
	if err != nil {
		log.Fatal(err)
	}
}
