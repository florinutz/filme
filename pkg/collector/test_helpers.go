package collector

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	gzBoltHttp "github.com/florinutz/gz-boltdb/http"

	"github.com/gocolly/colly"
)

const BucketName = "store"

// MockResponse parses the responses db and returns the mock response for a request
func MockResponse(request *http.Request, dataFile string) (response *colly.Response, err error) {
	httpResponse, err := gzBoltHttp.GetResponseFor(dataFile, BucketName, request, func(r1, r2 *http.Request) bool {
		return r1.URL.String() == r2.URL.String()
	})
	if err != nil {
		return
	}

	b, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}
	err = httpResponse.Body.Close()
	if err != nil {
		return nil, err
	}

	return &colly.Response{
		Request: &colly.Request{
			URL: httpResponse.Request.URL,
		},
		Body: b,
	}, nil
}

func UpdateTestData(reqs []*http.Request, outputPath string) []error {
	return gzBoltHttp.DumpResponses(reqs, outputPath, BucketName, nil, func(response *http.Response) error {
		fmt.Printf("* received %s\n", response.Request.URL.String())
		return nil
	})
}

func FatalIfErr(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}

func GenerateSimpleRequests(urls []string, tweakFn func(req *http.Request)) (reqs []*http.Request, err error) {
	for _, u := range urls {
		var req *http.Request
		req, err = http.NewRequest("GET", u, nil)
		if err != nil {
			err = fmt.Errorf("could not create a request for url '%s'", u)
			return
		}
		if tweakFn != nil {
			tweakFn(req)
		}
		reqs = append(reqs, req)
	}
	return
}

func GenerateRequestFromUrl(url string) *http.Request {
	reqs, err := GenerateSimpleRequests([]string{url}, func(req *http.Request) {
		req.Header.Set("Accept-Language", "en-US;q=0.8,es;q=0.5,fr;q=0.3")
	})
	if err != nil {
		panic("should not error here")
	}
	return reqs[0]
}
