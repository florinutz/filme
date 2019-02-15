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

type mockLoader struct {
	DataFile string
	Urls     []*Url
}

func NewMockLoader(dataFile string) *mockLoader {
	return &mockLoader{DataFile: dataFile}
}

type Url struct {
	Url     string `json:"url"`
	Content []byte `json:"content"`
}

// LoadFromFile reads all urls from data file
func (td *mockLoader) LoadFromFile() (err error) {
	jsonFile, err := os.Open(td.DataFile)
	if err != nil {
		return
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return
	}

	err = json.Unmarshal(byteValue, &td.Urls)
	if err != nil {
		return err
	}

	return
}

// Fetch loads new data
func (td *mockLoader) Fetch(wantedUrls []*url.URL, timeout time.Duration) error {
	c := make(chan Url, len(wantedUrls))

	for _, u := range wantedUrls {
		go fetch(*u, c)
	}

	for i := 0; i < len(wantedUrls); i++ {
		select {
		case block := <-c:
			td.Urls = append(td.Urls, &block)
			fmt.Printf("* loaded %s\n", block.Url)
		case <-time.After(timeout):
			return fmt.Errorf("timeout after %s", timeout)
		}
	}

	return nil
}

func (td *mockLoader) Save() error {
	b, err := json.Marshal(td.Urls)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(td.DataFile, b, 0644)
	if err != nil {
		return err
	}

	return nil
}

func fetch(url url.URL, blockChan chan<- Url) {
	urlStr := url.String()
	html, err := getSource(urlStr)
	if err != nil {
		log.Fatal(err)
	}
	result := make([]byte, base64.StdEncoding.EncodedLen(len(html)))
	base64.StdEncoding.Encode(result, html)
	blockChan <- Url{
		Url:     urlStr,
		Content: result,
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
