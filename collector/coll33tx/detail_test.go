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

func TestTorrent_fromResponse_Title(t *testing.T) {
	var torrent Torrent
	const pageLink = "https://1337x.to/torrent/3569899/House-Party-1990-WEBRip-720p-YTS-YIFY/"

	torrent.fromResponse(mockDetailPageResponse(pageLink), log.NewEntry(log.New()))

	if torrent.Title != "House Party (1990) [WEBRip] [1080p] [YTS] [YIFY]" {
		t.Fatal("wrong Title")
	}

	if torrent.FilmTitle != "House Party" {
		t.Fatal("wrong FilmTitle")
	}

	if torrent.FoundOn.String() != pageLink {
		t.Fatal("wrong PageLink")
	}

	if torrent.Magnet != "magnet:?xt=urn:btih:7F736E2E527ADF321D94092AA3DDA2C326EF4F31&dn=House+Party+%281990%29+%5BWEBRip%5D+%5B1080p%5D+%5BYTS%5D+%5BYIFY%5D&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969%2Fannounce&tr=udp%3A%2F%2F9.rarbg.com%3A2710%2Fannounce&tr=udp%3A%2F%2Fp4p.arenabg.com%3A1337&tr=udp%3A%2F%2Ftracker.internetwarriors.net%3A1337&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.zer0day.to%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969%2Fannounce&tr=udp%3A%2F%2Fcoppersurfer.tk%3A6969%2Fannounce" {
		t.Fatal("wrong Magnet")
	}

	if torrent.Category != "Movies" {
		t.Fatal("wrong Category")
	}

	if torrent.Type != "HD" {
		t.Fatal("wrong Type")
	}

	if torrent.Language != "English" {
		t.Fatal("wrong Language")
	}

	if torrent.TotalSize != "1.7 GB" {
		t.Fatal("wrong TotalSize")
	}
}
