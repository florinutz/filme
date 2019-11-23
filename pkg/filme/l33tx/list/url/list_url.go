package url

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/florinutz/filme/pkg/config/value/1337x/list/encoding"
	"github.com/florinutz/filme/pkg/config/value/1337x/list/search_category"
	"github.com/florinutz/filme/pkg/config/value/1337x/list/sort"
)

func GetListUrl(
	search string,
	sort sort.Value,
	category *search_category.SearchCategory,
	movieEncoding *encoding.ListEncoding,
) (u *url.URL, err error) {

	if search == "" { // must be a movies/encoding search
		if movieEncoding == nil {
			return nil, fmt.Errorf("you have to be looking for some encoding if there's no search query")
		}
		if category != nil {
			return nil, fmt.Errorf("category should be nil for an encoding listing")
		}

		return getEncodingUrl(*movieEncoding, sort), nil
	}

	return getSearchUrl(search, category, sort), nil
}

func getSearchUrl(search string, category *search_category.SearchCategory, sort sort.Value) *url.URL {
	const tpl = "https://1337x.to/sort-category-search/game%20of%20thrones%20s04e03/TV/time/desc/1/"
	u, _ := url.Parse(tpl)
	p := strings.Split(u.Path, "/")

	p[2] = search

	var (
		noCategory           = category == nil || category.TranslateToUrlParam() == ""
		indexForSortCriteria = 4
	)

	if noCategory {
		// change action
		p[1] = "sort-search"
		// remove category from url (index 3)
		p = append(p[:3], p[4:]...)
		indexForSortCriteria--
	} else {
		p[3] = category.TranslateToUrlParam()
	}

	p[indexForSortCriteria] = sort.Criteria.String()
	p[indexForSortCriteria+1] = sort.Order.String()

	u.Path = strings.Join(p, "/")

	return u
}

func getEncodingUrl(encoding encoding.ListEncoding, sort sort.Value) *url.URL {
	const tpl = "https://1337x.to/sort-sub/42/seeders/desc/1/"

	u, _ := url.Parse(tpl)
	p := strings.Split(u.Path, "/")

	p[2] = strconv.Itoa(int(encoding))
	p[3] = sort.Criteria.String()
	p[4] = sort.Order.String()

	u.Path = strings.Join(p, "/")

	return u
}
