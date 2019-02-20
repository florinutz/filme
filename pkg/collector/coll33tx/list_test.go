package coll33tx

import (
	"net/url"
	"testing"

	"github.com/florinutz/filme/pkg/collector/coll33tx/html/mockloader"
	"github.com/gocolly/colly"
)

const (
	TestPageDetail             = "https://1337x.to/torrent/3570061/House-Party-1990-WEBRip-1080p-YTS-YIFY/"
	TestPageList               = "https://1337x.to/popular-movies"
	TestPageListWithPagination = "https://1337x.to/search/romania/3/"
)

func mockResponse(pageUrl string) (*colly.Response, error) {
	u, err := url.Parse(pageUrl)
	if err != nil {
		return nil, err
	}

	loader := mockloader.NewMockLoader("html/data.json")

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

func Test_ListPage(t *testing.T) {
	responseNoPagination, err := mockResponse(TestPageList)
	fatalIfErr(err, t)
	docNoPagination, err := NewListPageDocument(responseNoPagination)
	fatalIfErr(err, t)
	responseWithPagination, err := mockResponse(TestPageListWithPagination)
	fatalIfErr(err, t)
	docWithPagination, err := NewListPageDocument(responseWithPagination)
	fatalIfErr(err, t)

	t.Run("pagination - response without pagination", func(t *testing.T) {
		_, err = docNoPagination.GetPagination()
		if err == nil {
			t.Error("this is supposed to return an error if the pagination is missing from a listing page")
		}
	})

	t.Run("pagination - response has pagination", func(t *testing.T) {
		p, err := docWithPagination.GetPagination()
		fatalIfErr(err, t)

		expected := 34
		if p.totalPages != expected {
			t.Errorf("wrong pages count (expected %d, got %d)", expected, p.totalPages)
		}
		expected = 3
		if p.currentPage != expected {
			t.Errorf("wrong current page (expected %d, got %d)", expected, p.currentPage)
		}
	})

	t.Run("items", func(t *testing.T) {
		items, err := docNoPagination.GetPageItems()
		fatalIfErr(err, t)

		if len(items) != 50 {
			// todo test individual items as general as possible, so the tests won't fail
			//  every time the list is refreshed with different items
			t.Fatalf("got %d items, expected %d", len(items), 50)
		}
	})
}

func fatalIfErr(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}
