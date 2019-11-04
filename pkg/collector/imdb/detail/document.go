package detail

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/florinutz/filme/pkg/collector"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

const (
	TestPageFilm             = "https://www.imdb.com/title/tt0497465/"
	TestPageSeriesUnfinished = "https://www.imdb.com/title/tt9561862/"
	TestPageSeriesFinished   = "https://www.imdb.com/title/tt0077000/"
)

// ImdbFilm represents the data found on a detail page
type ImdbFilm struct {
	Title string
	Year  int
}

type document struct {
	*goquery.Document
	log *log.Entry
}

func NewDocument(r *colly.Response, log *log.Entry) (*document, error) {
	d, err := collector.GetResponseDocument(r)
	if err != nil {
		return nil, err
	}
	return &document{d, log}, nil
}

func (doc *document) getTitle() (title string, err error) {
	selector := ".title_wrapper h1"
	selection := doc.Find(selector)
	if selection.Nodes == nil {
		err = fmt.Errorf("couldn't find the title: no element at selector '%s'", selector)
		return
	}
	selection.Children().Remove()
	title = strings.TrimSpace(selection.Text())

	return
}

func (doc *document) getYear() (year int, err error) {
	selector := "#titleYear a"
	selection := doc.Find(selector)
	if selection.Nodes == nil {
		err = fmt.Errorf("couldn't find the year: no element at selector '%s'", selector)
		return
	}

	return strconv.Atoi(selection.Text())
}

func (doc *document) GetData() (data *ImdbFilm) {
	var err error

	data = new(ImdbFilm)

	if data.Title, err = doc.getTitle(); err != nil {
		doc.log.WithError(err).Debug()
	}

	return
}
