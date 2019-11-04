package list

import (
	"testing"

	"github.com/florinutz/filme/pkg/collector"
)

const (
	DataFile = "../test-data"
)

func Test_ListPage(t *testing.T) {
	responseNoPagination, err := collector.MockResponse(
		collector.GenerateRequestFromUrl(TestPageList),
		DataFile,
	)
	collector.FatalIfErr(err, t)
	docNoPagination, err := NewDocument(responseNoPagination)
	collector.FatalIfErr(err, t)
	responseWithPagination, err := collector.MockResponse(
		collector.GenerateRequestFromUrl(TestPageListWithPagination),
		DataFile,
	)
	collector.FatalIfErr(err, t)
	docWithPagination, err := NewDocument(responseWithPagination)
	collector.FatalIfErr(err, t)

	t.Run("pagination - response without pagination", func(t *testing.T) {
		if docNoPagination.GetPagination() != nil {
			t.Error("document without pagination has pagination?")
		}
	})

	t.Run("pagination - response has pagination", func(t *testing.T) {
		p := docWithPagination.GetPagination()
		if p == nil {
			t.Error("pagination not found")
		}

		expected := 35
		if p.PagesCount != expected {
			t.Errorf("wrong pages count (expected %d, got %d)", expected, p.PagesCount)
		}
		expected = 3
		if p.Current != expected {
			t.Errorf("wrong current page (expected %d, got %d)", expected, p.Current)
		}
	})

	t.Run("items", func(t *testing.T) {
		lines := docNoPagination.GetLines()
		if lines == nil {
			t.Error("no lines found in list")
		}

		if len(lines) != 50 {
			// todo test individual items as general as possible, so the tests won't fail
			//  every time the list is refreshed with different items
			t.Fatalf("got %d items, expected %d", len(lines), 50)
		}
	})
}
