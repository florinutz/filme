package list

import (
	"reflect"
	"testing"
)

func Test_paging_getNextPages(t *testing.T) {
	type fields struct {
		filterLow    int
		filterHigh   int
		limitLow     int
		limitHigh    int
		itemsPerPage int
	}
	type args struct {
		wantedItems int
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantPages []int
	}{
		{
			name: "filterLow",
			fields: fields{
				filterLow:    3,
				filterHigh:   0,
				limitLow:     2,
				limitHigh:    5,
				itemsPerPage: 3,
			},
			args:      args{7},
			wantPages: []int{3, 4, 5},
		},
		{
			name: "filterHigh",
			fields: fields{
				filterLow:    0,
				filterHigh:   4,
				limitLow:     2,
				limitHigh:    5,
				itemsPerPage: 3,
			},
			args:      args{7},
			wantPages: []int{2, 3, 4},
		},
		{
			name: "full range",
			fields: fields{
				filterLow:    3,
				filterHigh:   5,
				limitLow:     1,
				limitHigh:    6,
				itemsPerPage: 3,
			},
			args:      args{10},
			wantPages: []int{3, 4, 5},
		},
		{
			name: "full range, restricted by items count",
			fields: fields{
				filterLow:    3,
				filterHigh:   5,
				limitLow:     1,
				limitHigh:    6,
				itemsPerPage: 3,
			},
			args:      args{2},
			wantPages: []int{3, 4, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Paging{
				filterLow:    tt.fields.filterLow,
				filterHigh:   tt.fields.filterHigh,
				limitLow:     tt.fields.limitLow,
				limitHigh:    tt.fields.limitHigh,
				itemsPerPage: tt.fields.itemsPerPage,
			}
			if gotPages := p.getNextPages(tt.args.wantedItems); !reflect.DeepEqual(gotPages, tt.wantPages) {
				t.Errorf("getNextPages() = %v, want %v", gotPages, tt.wantPages)
			}
		})
	}
}
