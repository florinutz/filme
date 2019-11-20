package l33tx_movies

import (
	"reflect"
	"testing"

	"github.com/florinutz/filme/pkg/config/value"
)

func TestSearchListUrl(t *testing.T) {
	var (
		searchCategoryTV  = value.SearchCategoryTV
		searchCategoryXXX = value.SearchCategoryXXX
		encodingUHD       = value.EncodingUHD
	)

	type args struct {
		search   string
		sort     value.LeetxListSortValue
		category *value.LeetxListSearchCategory
		encoding *value.LeetxListEncoding
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"TV time-desc",
			args{
				"game of thrones s04e03",
				value.LeetxListSortValue{Criteria: value.SortCriteriaTime, Order: value.SortOrderDesc},
				&searchCategoryTV,
				nil,
			},
			"https://1337x.to/sort-category-search/game%20of%20thrones%20s04e03/TV/time/desc/1/",
			false,
		},
		{
			"XXX size-asc",
			args{
				"game of thrones s04e03",
				value.LeetxListSortValue{Criteria: value.SortCriteriaSize, Order: value.SortOrderAsc},
				&searchCategoryXXX,
				nil,
			},
			"https://1337x.to/sort-category-search/game%20of%20thrones%20s04e03/XXX/size/asc/1/",
			false,
		},
		{
			"encoding UHD",
			args{
				"",
				value.LeetxListSortValue{Criteria: value.SortCriteriaSeeders, Order: value.SortOrderDesc},
				nil,
				&encodingUHD,
			},
			"https://1337x.to/sort-sub/76/seeders/desc/1/",
			false,
		},
		{
			"error for empty search and encoding",
			args{
				"",
				value.LeetxListSortValue{Criteria: value.SortCriteriaSize, Order: value.SortOrderAsc},
				&searchCategoryXXX,
				nil,
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUrl, err := GetListUrl(tt.args.search, tt.args.sort, tt.args.category, tt.args.encoding)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if gotUrl == nil {
				t.Fatalf("an url should be returned here")
				return
			}
			if !reflect.DeepEqual(gotUrl.String(), tt.want) {
				t.Errorf("gotUrl = %v, want %v", gotUrl, tt.want)
			}
		})
	}
}
