package coll33tx

import (
	"encoding/base64"
	"io/ioutil"
	"net/url"
	"testing"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

var DetailResponse *colly.Response

func init() {
	src, err := ioutil.ReadFile("html/detail")
	if err != nil {
		log.Fatal(err)
	}

	data, err := base64.StdEncoding.DecodeString(string(src))
	if err != nil {
		log.Fatal(err)
	}

	detailPageUrl, _ := url.Parse("https://1337x.to/torrent/3569899/House-Party-1990-WEBRip-720p-YTS-YIFY/")

	DetailResponse = &colly.Response{
		Body:    data,
		Request: &colly.Request{URL: detailPageUrl},
	}
}

func TestTorrent_fromResponse(t *testing.T) {
	type fields struct {
		Title           string
		FilmTitle       string
		FilmLink        *url.URL
		FilmCategories  []string
		FilmDescription string
		IMDB            *url.URL
		Magnet          string
		FoundOn         *url.URL
		Image           *url.URL
		Description     string
		Genres          []string
		Category        string
		Type            string
		Language        string
		TotalSize       string
		Downloads       int
		LastChecked     string
		DateUploaded    string
		Seeds           int
		Leeches         int
	}
	type args struct {
		r           *colly.Response
		responseLog *log.Entry
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "basic",
			args: struct {
				r           *colly.Response
				responseLog *log.Entry
			}{r: DetailResponse, responseLog: log.NewEntry(log.New())},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			torrent := &Torrent{
				Title:           tt.fields.Title,
				FilmTitle:       tt.fields.FilmTitle,
				FilmLink:        tt.fields.FilmLink,
				FilmCategories:  tt.fields.FilmCategories,
				FilmDescription: tt.fields.FilmDescription,
				IMDB:            tt.fields.IMDB,
				Magnet:          tt.fields.Magnet,
				FoundOn:         tt.fields.FoundOn,
				Image:           tt.fields.Image,
				Description:     tt.fields.Description,
				Genres:          tt.fields.Genres,
				Category:        tt.fields.Category,
				Type:            tt.fields.Type,
				Language:        tt.fields.Language,
				TotalSize:       tt.fields.TotalSize,
				Downloads:       tt.fields.Downloads,
				LastChecked:     tt.fields.LastChecked,
				DateUploaded:    tt.fields.DateUploaded,
				Seeds:           tt.fields.Seeds,
				Leeches:         tt.fields.Leeches,
			}

			torrent.fromResponse(tt.args.r, tt.args.responseLog)

			// todo add failure cases
		})
	}
}
