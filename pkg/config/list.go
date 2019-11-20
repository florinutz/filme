package config

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// https://www.google.com/search?q=imdb|rottentomatoes|metacritic|netflix|wikipedia+love+death+robots&num=20&lr=lang_en
func GetGoogleUrlFromArgs(args []string, numberOfResults int, language string) (string, error) {
	tplUrl := fmt.Sprintf("https://www.google.com/search?num=%d&lr=%s", numberOfResults, language) +
		"&q=imdb|rottentomatoes|metacritic|netflix|wikipedia+%s"

	var u string
	if u = GetUrl(args, `google\.com`, tplUrl, ""); u == "" {
		return "", errors.New("nothing to search for")
	}

	return u, nil
}

func GetUrl(args []string, domain, templateUrl, defaultUrl string) string {
	if len(args) == 0 {
		return defaultUrl
	}

	urlExtractor := func(what, domain string) string {
		u, err := url.Parse(what)
		if err != nil {
			return ""
		}
		if u.Hostname() == domain {
			return u.String()
		}
		return ""
	}

	fullString := strings.Join(args, " ")

	urlStr := urlExtractor(fullString, domain)
	if urlStr == "" {
		return fmt.Sprintf(
			templateUrl,
			url.QueryEscape(fullString),
		)
	}

	return urlStr
}
