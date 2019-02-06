package coll33tx

import (
	"encoding/base64"
	"io/ioutil"
	"net/url"
	"testing"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

func mockDetailPageResponse(pageUrl string) *colly.Response {
	src, err := ioutil.ReadFile("html/detail")
	if err != nil {
		log.Fatal(err)
	}

	data, err := base64.StdEncoding.DecodeString(string(src))
	if err != nil {
		log.Fatal(err)
	}
	u, _ := url.Parse(pageUrl)

	return &colly.Response{
		Body:    data,
		Request: &colly.Request{URL: u},
	}
}

func TestTorrent_fromResponse(t *testing.T) {
	var torrent Torrent
	response := mockDetailPageResponse("https://1337x.to/torrent/3569899/House-Party-1990-WEBRip-720p-YTS-YIFY/")
	torrent.fromResponse(response, nil)
}
