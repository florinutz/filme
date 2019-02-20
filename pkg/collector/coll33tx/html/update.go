package main

import (
	"log"
	"net/url"
	"time"

	"github.com/florinutz/filme/pkg/collector/coll33tx/html/mockloader"
)

const dataFile = "data.json"

// todo make this a separate package
func main() {
	loader := mockloader.NewMockLoader(dataFile)

	url1, _ := url.Parse("https://1337x.to/torrent/3570061/House-Party-1990-WEBRip-1080p-YTS-YIFY/")
	url2, _ := url.Parse("https://1337x.to/search/romania/3/")
	url3, _ := url.Parse("https://1337x.to/popular-movies")
	urls := []*url.URL{url1, url2, url3}

	timeout := 10 * time.Second

	err := loader.Fetch(urls, timeout)
	if err != nil {
		log.Fatal(err)
	}

	err = loader.Save()
	if err != nil {
		log.Fatal(err)
	}

	// content, err := loader.GetUrlContent(url1)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("%s", content)
}
