package l33tx_movies

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/florinutz/filme/pkg/config/value"
)

func GetListUrl(
	search string,
	sort value.LeetxListSortValue,
	category *value.LeetxListSearchCategory,
	movieEncoding *value.LeetxListEncoding,
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

	if category == nil {
		return nil, fmt.Errorf("you need a category here")
	}

	return getSearchUrl(search, *category, sort), nil
}

func getSearchUrl(search string, category value.LeetxListSearchCategory, sort value.LeetxListSortValue) *url.URL {
	const tpl = "https://1337x.to/sort-category-search/game%20of%20thrones%20s04e03/TV/time/desc/1/"
	u, _ := url.Parse(tpl)
	p := strings.Split(u.Path, "/")

	p[2] = search
	p[3] = category.TranslateToUrlParam()
	p[4] = sort.Criteria.String()
	p[5] = sort.Order.String()

	u.Path = strings.Join(p, "/")

	return u
}

func getEncodingUrl(encoding value.LeetxListEncoding, sort value.LeetxListSortValue) *url.URL {
	const tpl = "https://1337x.to/sort-sub/42/seeders/desc/1/"

	u, _ := url.Parse(tpl)
	p := strings.Split(u.Path, "/")

	p[2] = strconv.Itoa(int(encoding))
	p[3] = sort.Criteria.String()
	p[4] = sort.Order.String()

	u.Path = strings.Join(p, "/")

	return u
}
