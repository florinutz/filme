package collector

import (
	"errors"

	"github.com/gocolly/colly"
)

var (
	ErrUnknownResponse          = errors.New("unknown response")
	ErrMultipleResponseHandlers = errors.New("you have multiple crawlers that can handle this response")
)

//ResponseMatcher sees interesting data in a colly response
type CrawlerResponseHandler interface {
	CrawlerGenerator
	CanHandleResponse(*colly.Response) bool
}

type CrawlerGenerator interface {
	Create(response *colly.Response) (*colly.Collector, error)
}

type CrawlerFactory struct {
	Factories []CrawlerResponseHandler
}

func (cf CrawlerFactory) Create(response *colly.Response) (*colly.Collector, error) {
	var matches []CrawlerResponseHandler

	for _, responseHandler := range cf.Factories {
		if responseHandler.CanHandleResponse(response) {
			matches = append(matches, responseHandler)
		}
	}

	if len(matches) > 1 {
		return nil, ErrMultipleResponseHandlers
	}

	if matches == nil {
		return nil, ErrUnknownResponse
	}

	return matches[0].Create(response)
}
