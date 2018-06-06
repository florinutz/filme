package coll33tx

import (
	"net/url"
	"strings"
	"sync"

	"bytes"

	"errors"

	"strconv"

	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

const BlankImage = "32e6dd6abe806d43e9453adf3d310851.jpg"

var imdbRe = regexp.MustCompile(`(?m)https?://(www\.)?imdb.com/title/tt(\d)+`)

// L33tTorrent represents the data onItemFound on a torrent details page
type L33tTorrent struct {
	Title           string
	FilmTitle       string
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
	Downloads       int
	LastChecked     string
	DateUploaded    string
	Seeds           int
	Leeches         int
}

// TorrentFoundCallback is the type the callback func that's be called when a torrent was onItemFound
type TorrentFoundCallback func(torrent L33tTorrent)

// detailsCollector is a wrapper around the colly collector + page data
type detailsCollector struct {
	*colly.Collector
	found   TorrentFoundCallback
	torrent L33tTorrent // this will be filled in the events
}

var (
	mutex *sync.Mutex
)

func init() {
	mutex = &sync.Mutex{}
}

func NewDetailsPageCollector(found TorrentFoundCallback, options ...func(*colly.Collector)) *detailsCollector {
	col := detailsCollector{
		Collector: getCollyCollector(options...),
		found:     found,
		torrent:   L33tTorrent{},
	}

	col.Collector.OnResponse(col.OnResponse)
	col.Collector.OnScraped(col.OnScraped)

	return &col
}

// Magnet gets the magnet link from the details page
func (c *detailsCollector) Magnet(e *colly.HTMLElement) {
	if !strings.HasPrefix(e.Attr("href"), "magnet") {
		return
	}
	c.torrent.Magnet = e.Attr("href")
}

func (c *detailsCollector) OnResponse(r *colly.Response) {
	c.torrent.fromResponse(r)
}

// OnScraped assembles and collects the Torrent struct at the end
func (c *detailsCollector) OnScraped(r *colly.Response) {
	c.found(c.torrent)
}

func (torrent *L33tTorrent) fromResponse(r *colly.Response) (errs []error) {
	mutex.Lock()
	defer mutex.Unlock()

	torrent.FoundOn = r.Request.URL

	pageLog := log.WithField("url", torrent.FoundOn.String())

	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(r.Body))
	if err != nil {
		return []error{err}
	}

	if title := doc.Find(".box-info-heading h1"); title.Nodes == nil {
		errs = append(errs, errors.New("missing title element"))
	} else {
		torrent.Title = title.Text()
	}

	if links := doc.Find("a[href]").FilterFunction(func(_ int, s *goquery.Selection) bool {
		href, _ := s.Attr("href")
		return strings.HasPrefix(href, "magnet:?")
	}); links.Nodes == nil {
		errs = append(errs, errors.New("missing magnet link element"))
	} else {
		torrent.Magnet, _ = links.Attr("href")
	}

	doc.Find("ul.list li").
		FilterFunction(func(_ int, s *goquery.Selection) bool {
			hasExactlyTwoChildren := s.Children().Length() == 2
			firstChildIsStrong := s.Children().Eq(0).Is("strong")
			secondChildIsSpan := s.Children().Eq(1).Is("span")

			return hasExactlyTwoChildren && firstChildIsStrong && secondChildIsSpan
		}).
		Each(func(_ int, s *goquery.Selection) {
			label := s.Children().Eq(0).Text()
			value := s.Children().Eq(1).Text()

			switch label {
			case "Category":
				torrent.Category = value
			case "Type":
				torrent.Type = value
			case "Language":
				torrent.Language = value
			case "Total size":
				torrent.TotalSize = value
			case "Downloads":
				torrent.Downloads, _ = strconv.Atoi(value)
			case "Date uploaded":
				torrent.DateUploaded = value
			case "Seeders":
				torrent.Seeds, _ = strconv.Atoi(value)
			case "Leechers":
				torrent.Leeches, _ = strconv.Atoi(value)
			}
		})

	if img := doc.Find(".torrent-detail .torrent-image img"); img.Nodes == nil {
		errs = append(errs, errors.New("no image element"))
	} else {
		if src, exists := img.Attr("src"); !exists {
			errs = append(errs, errors.New("image element has no src"))
		} else {
			if strings.HasSuffix(src, BlankImage) {
				pageLog.Debug("default image")
			} else {
				if torrent.Image, err = url.Parse(src); err != nil {
					errs = append(errs, errors.New("invalid image url"))
				} else {
					if strings.HasPrefix(torrent.Image.String(), "//") {
						torrent.Image, _ = url.Parse("http://" + torrent.Image.String()[2:])
					}
				}
			}
		}
	}

	if filmTitle := doc.Find(".torrent-detail-info h3 a"); filmTitle.Nodes != nil {
		torrent.FilmTitle = filmTitle.Text()
		link, _ := filmTitle.Attr("href")
		if torrent.FilmLink, err = url.Parse(r.Request.AbsoluteURL(link)); err != nil {
			errs = append(errs, errors.New("invalid normalized link"))
		}
	} else {
		errs = append(errs, errors.New("no normalized title element"))
	}

	if filmCategories := doc.Find(".torrent-category span"); filmCategories.Nodes != nil {
		filmCategories.Each(func(_ int, s *goquery.Selection) {
			torrent.FilmCategories = append(torrent.FilmCategories, s.Text())
		})
	} else {
		errs = append(errs, errors.New("no film categories"))
	}

	if filmDescription := doc.Find(".torrent-detail-info p"); filmDescription != nil {
		torrent.FilmDescription = filmDescription.Text()
	} else {
		errs = append(errs, errors.New("no film description"))
	}

	if matches := imdbRe.FindAllString(string(r.Body), -1); matches != nil {
		torrent.IMDB, _ = url.Parse(matches[0])
	}

	return
}
