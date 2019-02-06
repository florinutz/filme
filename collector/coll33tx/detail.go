package coll33tx

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"bytes"

	"strconv"

	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

const BlankImage = "32e6dd6abe806d43e9453adf3d310851.jpg"

var imdbRe = regexp.MustCompile(`(?m)https?://(www\.)?imdb.com/title/tt(\d)+`)

// Torrent represents the data onItemFound on a Torrent details page
type Torrent struct {
	Title           string
	FilmCleanTitle  string
	FilmLink        *url.URL
	FilmCategories  []string
	FilmDescription string
	IMDB            *url.URL
	Magnet          string
	FoundOn         *url.URL
	Image           *url.URL
	Description     string
	Genres          []string
	Category        string
	Type            string
	Language        string
	TotalSize       string
	UploadedBy      string
	Downloads       int
	LastChecked     string
	DateUploaded    string
	Seeders         int
	Leechers        int
}

// TorrentFoundCallback is the type the callback func that's be called when a Torrent was onItemFound
type TorrentFoundCallback func(torrent Torrent)

// DetailsCollector is a wrapper around the colly collector + page data
type DetailsCollector struct {
	*colly.Collector
	found   TorrentFoundCallback
	Torrent Torrent // this will be filled in the events
	Log     *log.Entry
}

var (
	mutex *sync.Mutex
)

func init() {
	mutex = &sync.Mutex{}
}

func NewDetailsPageCollector(found TorrentFoundCallback, log *log.Entry, options ...func(*colly.Collector)) *DetailsCollector {
	col := DetailsCollector{
		Collector: initCollector(log, options...),
		found:     found,
		Torrent:   Torrent{},
		Log:       log,
	}

	col.Collector.OnResponse(col.OnResponse)
	col.Collector.OnScraped(col.OnScraped)

	return &col
}

// Magnet gets the magnet link from the details page
func (dc *DetailsCollector) Magnet(e *colly.HTMLElement) {
	if !strings.HasPrefix(e.Attr("href"), "magnet") {
		return
	}
	dc.Torrent.Magnet = e.Attr("href")
}

func (dc *DetailsCollector) OnResponse(r *colly.Response) {
	dc.Torrent.fromResponse(r, dc.Log)
}

// OnScraped assembles and collects the Torrent struct at the end
func (dc *DetailsCollector) OnScraped(r *colly.Response) {
	dc.found(dc.Torrent)
}

func getDetailsPageDoc(r *colly.Response) (doc *goquery.Document, err error) {
	return goquery.NewDocumentFromReader(bytes.NewBuffer(r.Body))
}

func getDetailsPageTitle(doc *goquery.Document) (title string, err error) {
	selector := ".box-info-heading h1"
	selection := doc.Find(selector)
	if selection.Nodes == nil {
		err = fmt.Errorf("couldn't find the title: no element at selector %s", selector)
		return
	}
	title = strings.TrimSpace(selection.Text())
	return
}

func getDetailsPageMagnet(doc *goquery.Document) (magnet string, err error) {
	filterByMagnetSchemeFunc := func(_ int, s *goquery.Selection) bool {
		href, _ := s.Attr("href")
		return strings.HasPrefix(href, "magnet:?")
	}

	allLinks := doc.Find("a[href]")

	selection := allLinks.FilterFunction(filterByMagnetSchemeFunc)
	if selection.Nodes == nil {
		err = errors.New("missing magnet link")
		return
	}

	magnet, _ = selection.Attr("href")

	return magnet, nil
}

func getDetailsPageBox(doc *goquery.Document) (result *struct {
	Category     string
	Type         string
	Language     string
	TotalSize    string
	UploadedBy   string
	Downloads    int
	LastChecked  string
	DateUploaded string
	Seeders      int
	Leechers     int
}, err error) {
	infoItemFilterFunc := func(_ int, s *goquery.Selection) bool {
		hasExactlyTwoChildren := s.Children().Length() == 2
		firstChildIsStrong := s.Children().Eq(0).Is("strong")
		secondChildIsSpan := s.Children().Eq(1).Is("span")

		return hasExactlyTwoChildren && firstChildIsStrong && secondChildIsSpan
	}
	items := doc.Find("ul.list li").FilterFunction(infoItemFilterFunc)
	if items.Nodes == nil {
		return nil, errors.New("no box items found")
	}
	items.Each(func(_ int, s *goquery.Selection) {
		label := s.Children().Eq(0).Text()
		value := s.Children().Eq(1).Text()
		switch label {
		case "Category":
			result.Category = value
		case "Type":
			result.Type = value
		case "Language":
			result.Language = value
		case "Total size":
			result.TotalSize = value
		case "Uploaded By":
			result.UploadedBy, _ = s.Children().Eq(1).Find("a").Attr("href")
		case "Downloads":
			result.Downloads, _ = strconv.Atoi(value)
		case "Date uploaded":
			result.DateUploaded = value
		case "Seeders":
			result.Seeders, _ = strconv.Atoi(value)
		case "Leechers":
			result.Leechers, _ = strconv.Atoi(value)
		}
	})

	return
}

func getDetailsPageImage(doc *goquery.Document) (image *url.URL, err error) {
	const imgSelector = ".torrent-detail .torrent-image img"
	selection := doc.Find(imgSelector)
	if selection.Nodes == nil {
		return nil, fmt.Errorf("missing image element at selector '%s'", imgSelector)
	}
	src, exists := selection.Attr("src")
	if !exists {
		return nil, errors.New("image has no src attribute")
	}
	if strings.HasSuffix(src, BlankImage) {
		return nil, errors.New("image is the default one")
	}

	return url.Parse(src)
}

func getDetailsPageFilm(doc *goquery.Document, request *colly.Request) (
	film *struct {
		title string
		url   *url.URL
	},
	err error,
) {
	const selector = ".torrent-detail-info h3 a"
	selection := doc.Find(selector)
	if selection.Nodes == nil {
		return nil, fmt.Errorf("missing film info (title/link) at selector '%s'", selector)
	}
	film.title = selection.Text()

	href, exists := selection.Attr("href")
	if !exists {
		return film, errors.New("film title link has no href")
	}

	film.url, err = url.Parse(request.AbsoluteURL(href))

	return
}

func (torrent *Torrent) fromResponse(r *colly.Response, responseLog *log.Entry) {
	mutex.Lock()
	defer mutex.Unlock()

	torrent.FoundOn = r.Request.URL

	doc, err := getDetailsPageDoc(r)
	if err != nil {
		if responseLog != nil {
			responseLog.WithError(err).Fatal("couldn't init doc")
		}
	}

	if torrent.Title, err = getDetailsPageTitle(doc); err != nil && responseLog != nil {
		responseLog.WithError(err).Debug("missing title element")
	}

	if torrent.Magnet, err = getDetailsPageMagnet(doc); err != nil && responseLog != nil {
		responseLog.WithError(err).Debug("missing magnet")
	}

	if box, err := getDetailsPageBox(doc); err == nil {
		torrent.Category = box.Category
		torrent.Type = box.Type
		torrent.Language = box.Language
		torrent.TotalSize = box.TotalSize
		torrent.UploadedBy = box.UploadedBy
		torrent.Downloads = box.Downloads
		torrent.LastChecked = box.LastChecked
		torrent.DateUploaded = box.DateUploaded
		torrent.Seeders = box.Seeders
		torrent.Leechers = box.Leechers
	} else if responseLog != nil {
		responseLog.WithError(err).Warn()
	}

	if torrent.Image, err = getDetailsPageImage(doc); err != nil && responseLog != nil {
		responseLog.WithError(err).Debug()
	}

	film, err := getDetailsPageFilm(doc, r.Request)
	if err != nil && responseLog != nil {
		responseLog.WithError(err).Debug()
	}
	torrent.FilmCleanTitle = film.title
	torrent.FilmLink = film.url

	if filmCategories := doc.Find(".torrent-category span"); filmCategories.Nodes != nil {
		filmCategories.Each(func(_ int, s *goquery.Selection) {
			torrent.FilmCategories = append(torrent.FilmCategories, s.Text())
		})
	} else {
		if responseLog != nil {
			responseLog.Debug("no film categories")
		}
	}

	if filmDescription := doc.Find(".torrent-detail-info p"); filmDescription != nil {
		torrent.FilmDescription = filmDescription.Text()
	} else {
		if responseLog != nil {
			responseLog.Debug("no film description")
		}
	}

	if matches := imdbRe.FindAllString(string(r.Body), -1); matches != nil {
		torrent.IMDB, _ = url.Parse(matches[0])
	}

	return
}

func (dc *DetailsCollector) CanHandleResponse(r *colly.Response) bool {
	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(r.Body))
	if err != nil {
		log.WithField("response", r).Warn("error while checking if ListCollector can handle a response")
		return false
	}
	if title := doc.Find(".box-info-heading h1"); title.Nodes == nil {
		return false
	}
	if img := doc.Find(".torrent-detail .torrent-image img"); img.Nodes == nil {
		return false
	}
	return true
}
