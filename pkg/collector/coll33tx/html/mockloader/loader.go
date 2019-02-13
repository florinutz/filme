package mockloader

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type block struct {
	Url string `json:"url"`
	B64 []byte `json:"content"`
}

type Data struct {
	Detail                block `json:"detail"`
	ListWithPagination    block `json:"list_pagination"`
	ListWithoutPagination block `json:"list_nopagination"`
}

// Read Data from data file
func Read(dataFile string) (data *Data, err error) {
	jsonFile, err := os.Open(dataFile)
	if err != nil {
		return
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		return nil, err
	}

	return
}

// Write Data to file
func Write(dataFile string) error {
	c1 := make(chan block, 1)
	c2 := make(chan block, 1)
	c3 := make(chan block, 1)

	go fetch("https://1337x.to/torrent/3570061/House-Party-1990-WEBRip-1080p-YTS-YIFY/", c1)
	go fetch("https://1337x.to/search/romania/1/", c2)
	go fetch("https://1337x.to/popular-movies", c3)

	data := &Data{}

	for i := 0; i < 3; i++ {
		select {
		case data.Detail = <-c1:
			fmt.Println("loaded detail page")
		case data.ListWithPagination = <-c2:
			fmt.Println("loaded list page with pagination")
		case data.ListWithoutPagination = <-c3:
			fmt.Println("loaded list page without pagination")
		case <-time.After(10 * time.Second):
			log.Fatal("timeout after 10 seconds")
		}
	}

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dataFile, b, 0644)
	if err != nil {
		return err
	}

	return nil
}

func fetch(url string, blockChan chan<- block) {
	html, err := getSource(url)
	if err != nil {
		log.Fatal(err)
	}
	result := make([]byte, base64.StdEncoding.EncodedLen(len(html)))
	base64.StdEncoding.Encode(result, html)
	blockChan <- block{
		Url: url,
		B64: result,
	}
}

func getSource(visitUrl string) ([]byte, error) {
	_, err := url.Parse(visitUrl)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(visitUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return html, nil
}
