package list

import (
	"testing"

	"github.com/florinutz/filme/pkg/collector/gz_http"
)

const (
	DataFile = "../../../collector/coll33tx/test-data"
)

func Test_ListPage(t *testing.T) {
	responseNoPagination, err := gz_http.MockResponse(
		gz_http.GenerateRequestFromUrl(TestPageList),
		DataFile,
	)
	gz_http.FatalIfErr(err, t)
	docNoPagination, err := NewDocument(responseNoPagination)
	gz_http.FatalIfErr(err, t)
	responseWithPagination, err := gz_http.MockResponse(
		gz_http.GenerateRequestFromUrl(TestPageListWithPagination),
		DataFile,
	)
	gz_http.FatalIfErr(err, t)
	docWithPagination, err := NewDocument(responseWithPagination)
	gz_http.FatalIfErr(err, t)

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

		expected := 39
		if p.PagesCount != expected {
			t.Errorf("wrong pages count (expected %d, got %d)", expected, p.PagesCount)
		}
		expected = 3
		if p.Current != expected {
			t.Errorf("wrong current page (expected %d, got %d)", expected, p.Current)
		}
	})

	t.Run("items", func(t *testing.T) {
		lines, err := docNoPagination.GetLines()
		if err != nil {
			t.Errorf("lines fetching error: %s", err.Error())
		}
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
