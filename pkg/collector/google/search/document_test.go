// +build google

package search

import (
	"log"
	"testing"

	"github.com/florinutz/filme/pkg/collector"
	"github.com/gocolly/colly"
)

const (
	TestDataFile = "../test-data"
)

var (
	r     *colly.Response
	doc   *document
	items map[int]BaseItem
)

func init() {
	var err error

	r, err = collector.MockResponse(collector.GenerateRequestFromUrl(TestSearch), TestDataFile)
	if err != nil {
		log.Fatal(err)
	}
	doc, err = NewDocument(r, nil)
	if err != nil {
		log.Fatal(err)
	}
	items, err = doc.GetItems()
	if err != nil {
		log.Fatal(err)
	}
}

func TestSearchPage_Items_Count(t *testing.T) {
	if len(items) == 0 {
		t.Fatal("no items retrieved")
	}
	expected := 31
	got := len(items)
	if expected != got {
		t.Errorf("expected count %d, got %d", expected, got)
	}
}

func TestSearchPage_Item_Title(t *testing.T) {
	item := getItemOrSkip(t)
	expected := "Love, Death & Robots - Wikipedia"
	got := item.Title()
	if expected != got {
		t.Errorf("expected title '%s', got '%s'", expected, got)
	}
}

func TestSearchPage_Item_Url(t *testing.T) {
	item := getItemOrSkip(t)
	expected := "https://en.wikipedia.org/wiki/Love,_Death_%26_Robots"
	got := item.Url().String()
	if expected != got {
		t.Errorf("expected url '%s', got '%s'", expected, got)
	}
}

func TestSearchPage_Item_Description(t *testing.T) {
	item := getItemOrSkip(t)
	expected := ``
	got := item.Description()
	if expected != got {
		t.Errorf("expected url '%s', got '%s'", expected, got)
	}
}

func getItemOrSkip(t *testing.T) BaseItem {
	if len(items) == 0 {
		t.Skip("no items")
	}
	item := items[0]
	return item
}
