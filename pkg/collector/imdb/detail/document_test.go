package detail

import (
	"testing"

	"github.com/florinutz/filme/pkg/collector/gz_http"
)

const (
	DataFile = "../test-data"
)

func TestDetailPage_document_data(t *testing.T) {
	r, err := gz_http.MockResponse(
		gz_http.GenerateRequestFromUrl(TestPageFilm),
		DataFile,
	)
	gz_http.FatalIfErr(err, t)

	doc, err := NewDocument(r, nil)
	gz_http.FatalIfErr(err, t)

	data := doc.GetData()

	t.Run("Title", func(t *testing.T) {
		expected := "Vicky Cristina Barcelona"
		got := data.Title
		if expected != got {
			t.Errorf("expected Title '%s', got '%s'", expected, got)
		}
	})

	t.Run("Year", func(t *testing.T) {
		t.SkipNow()
		expected := 2019
		got := data.Year
		if expected != got {
			t.Errorf("expected year '%d', got '%d'", expected, got)
		}
	})
}
