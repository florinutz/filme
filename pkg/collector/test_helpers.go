package collector

import (
	"net/url"
	"testing"

	"github.com/florinutz/filme/pkg/collector/coll33tx/html/mockloader"
	"github.com/gocolly/colly"
)

func MockResponse(pageUrl string, dataFile string) (*colly.Response, error) {
	u, err := url.Parse(pageUrl)
	if err != nil {
		return nil, err
	}

	loader := mockloader.NewMockLoader(dataFile)

	err = loader.LoadFromFile()
	if err != nil {
		return nil, err
	}

	content, err := loader.GetUrlContent(u)
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
