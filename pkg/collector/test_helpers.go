package collector

import (
	"net/url"
	"testing"

	"github.com/florinutz/filme/pkg/collector/url_mock"

	"github.com/gocolly/colly"
)

func MockResponse(pageUrl string, dataFile string) (*colly.Response, error) {
	u, err := url.Parse(pageUrl)
	if err != nil {
		return nil, err
	}

	err = url_mock.Load(dataFile)
	if err != nil {
		return nil, err
	}

	content, err := url_mock.GetUrlContent(u.String())
	if err != nil {
		return nil, err
	}

	return &colly.Response{
		Body:    content,
		Request: &colly.Request{URL: u},
	}, nil
}

func FatalIfErr(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}
