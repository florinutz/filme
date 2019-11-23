package url

import (
	"reflect"
	"testing"

	"github.com/florinutz/filme/pkg/config/value/1337x/list/encoding"
	"github.com/florinutz/filme/pkg/config/value/1337x/list/search_category"
	"github.com/florinutz/filme/pkg/config/value/1337x/list/sort"
)

func TestSearchListUrl(t *testing.T) {
	var (
		searchCategoryAll = search_category.SearchCategoryAll
		searchCategoryTV  = search_category.SearchCategoryTV
		searchCategoryXXX = search_category.SearchCategoryXXX
		encodingUHD       = encoding.EncUHD
	)

	type args struct {
		search   string
		sort     sort.Value
		category *search_category.SearchCategory
		encoding *encoding.ListEncoding
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
				sort.Value{Criteria: sort.CriteriaTime, Order: sort.OrderDesc},
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
				sort.Value{Criteria: sort.CriteriaSize, Order: sort.OrderAsc},
				&searchCategoryXXX,
				nil,
			},
			"https://1337x.to/sort-category-search/game%20of%20thrones%20s04e03/XXX/size/asc/1/",
			false,
		},
		{
			"ALL with category",
			args{
				"game of thrones s04e03",
				sort.Value{Criteria: sort.CriteriaTime, Order: sort.OrderDesc},
				&searchCategoryAll,
				nil,
			},
			"https://1337x.to/sort-search/game%20of%20thrones%20s04e03/time/desc/1/",
			false,
		},
		{
			"ALL with nil category",
			args{
				"game of thrones s04e03",
				sort.Value{Criteria: sort.CriteriaTime, Order: sort.OrderDesc},
				nil,
				nil,
			},
			"https://1337x.to/sort-search/game%20of%20thrones%20s04e03/time/desc/1/",
			false,
		},
		{
			"encoding UHD",
			args{
				"",
				sort.Value{Criteria: sort.CriteriaSeeders, Order: sort.OrderDesc},
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
				sort.Value{Criteria: sort.CriteriaSize, Order: sort.OrderAsc},
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
