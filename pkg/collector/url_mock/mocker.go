package url_mock

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

var Urls []*Url

// Url - not using a map here so I can save this as json
type Url struct {
	Url     string `json:"url"`
	Content []byte `json:"content"`
}

// LoadFromFile reads all urls from data file
func Load(filePath string) (err error) {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return
	}

	err = json.Unmarshal(byteValue, &Urls)
	if err != nil {
		return err
	}

	return
}

// Fetch loads new data
func Fetch(wantedUrls []*url.URL, timeout time.Duration) error {
	c := make(chan Url, len(wantedUrls))

	for _, u := range wantedUrls {
		go fetch(*u, c)
	}

	for i := 0; i < len(wantedUrls); i++ {
		select {
		case block := <-c:
			SetUrlContent(block.Url, block.Content)
			fmt.Printf("* loaded %s\n", block.Url)
		case <-time.After(timeout):
			return fmt.Errorf("timeout after %s", timeout)
		}
	}

	return nil
}

func Persist(filePath string) error {
	b, err := json.Marshal(Urls)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, b, 0644)
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

func decodeB64(b64 []byte) (decoded []byte, err error) {
	decoded = make([]byte, base64.StdEncoding.DecodedLen(len(b64)))
	_, err = base64.StdEncoding.Decode(decoded, b64)
	if err != nil {
		return nil, errors.New("couldn't decode base64")
	}
	return decoded, nil
}

func GetUrlContent(u string) ([]byte, error) {
	if Urls == nil {
		return nil, errors.New("loader has no urls")
	}

	for _, loaderUrl := range Urls {
		if u == loaderUrl.Url {
			return decodeB64(loaderUrl.Content)
		}
	}

	return nil, fmt.Errorf("url '%s' not found", u)
}

// SetUrlContent adds a new url
func SetUrlContent(u string, content []byte) {
	for i, existing := range Urls {
		if u == existing.Url {
			Urls[i].Content = content
			return
		}
	}
	Urls = append(Urls, &Url{Url: u, Content: content})
}
