package detail

import (
	"testing"

	"github.com/florinutz/filme/pkg/collector"
)

const (
	DataFile = "../test-data"
)

func TestDetailPage_document_data(t *testing.T) {
	r, err := collector.MockResponse(
		collector.GenerateRequestFromUrl(TestPageFilm),
		DataFile,
	)
	collector.FatalIfErr(err, t)

	doc, err := NewDocument(r, nil)
	collector.FatalIfErr(err, t)

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
