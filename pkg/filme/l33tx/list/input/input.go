package input

import (
	"net/url"

	"github.com/florinutz/filme/pkg/config/value/1337x/list/encoding"
	"github.com/florinutz/filme/pkg/config/value/1337x/list/search_category"
	"github.com/florinutz/filme/pkg/config/value/1337x/list/sort"
	listurl "github.com/florinutz/filme/pkg/filme/l33tx/list/url"
)

// todo validation for the 2 cases from inside GetListUrl
type ListingInput struct {
	Search   string
	Url      *url.URL
	Category *search_category.SearchCategory
	Encoding *encoding.ListEncoding
	Sort     sort.Value
}

func (i ListingInput) GetStartUrl() (*url.URL, error) {
	return listurl.GetListUrl(
		i.Search,
		sort.Value{Criteria: sort.CriteriaSeeders, Order: sort.OrderDesc},
		i.Category,
		i.Encoding,
	)
}
