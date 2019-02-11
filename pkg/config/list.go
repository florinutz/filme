package config

import (
	"fmt"
	"net/url"
	"strings"
)

func GetListUrlFromArgs(args []string) (listUrl string) {
	if len(args) > 0 {
		firstArg := args[0]
		if strings.Contains(firstArg, "//1337x.to") {
			listUrl = firstArg
		} else { // search
			params := url.Values{"search": {strings.Join(args, " ")}}
			listUrl = fmt.Sprintf("https://1337x.to/srch?%s", params.Encode())
		}
	} else {
		listUrl = "https://1337x.to/popular-movies"
	}

	return
}
