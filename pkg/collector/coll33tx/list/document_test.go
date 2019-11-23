package list

import (
	"math"
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

func Test_ToBytes(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    uint64
		wantErr bool
	}{
		{
			name:    "simple",
			input:   "7096kb",
			want:    uint64(7096 * 1024),
			wantErr: false,
		},
		{
			name:    "multiple with space",
			input:   "12 GB",
			want:    uint64(12 * 1024 * 1024 * 1024),
			wantErr: false,
		},
		{
			name:    "decimals with dot",
			input:   "12.4 kb",
			want:    uint64(math.Floor(12.4 * 1024)),
			wantErr: false,
		},
		{
			name:    "decimals with comma",
			input:   "17,3 MiB",
			want:    uint64(math.Floor(17.3 * 1024 * 1024)),
			wantErr: false,
		},
		{
			name:    "bytes",
			input:   "17,3B",
			want:    uint64(math.Floor(17.3)),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToBytes(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToBytes() got = %v, want %v", got, tt.want)
			}
		})
	}
}
