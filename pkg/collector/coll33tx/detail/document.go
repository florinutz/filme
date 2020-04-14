package detail

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/florinutz/filme/pkg/collector"
	"github.com/florinutz/filme/pkg/collector/coll33tx"
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

const (
	BlankImage     = "32e6dd6abe806d43e9453adf3d310851.jpg"
	TestPageDetail = "https://1337x.to/torrent/3570061/House-Party-1990-WEBRip-1080p-YTS-YIFY/"
)

var imdbRE = regexp.MustCompile(`(?m)https?://(www\.)?imdb.com/title/tt(\d)+`)
var yearRE = regexp.MustCompile(`20\d{2}`)
var idRE = regexp.MustCompile(`/torrent/(\d+)/`)

type document struct {
	*goquery.Document
	log log.Entry
}

func NewDocument(r *colly.Response, log log.Entry) (*document, error) {
	d, err := collector.GetResponseDocument(r)
	if err != nil {
		return nil, err
	}
	return &document{d, log}, nil
}

func (doc *document) getDetailsPageTitle() (title string, err error) {
	selector := ".box-info-heading h1"
	selection := doc.Find(selector)
	if selection.Nodes == nil {
		err = fmt.Errorf("couldn't find the title: no element at selector %s", selector)
		return
	}
	title = strings.TrimSpace(selection.Text())
	return
}

func (doc *document) getDetailsPageMagnet() (magnet string, err error) {
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

func (doc *document) getDetailsPageBox() (result struct {
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
		err = errors.New("no box items found")
		return
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
			if val, hasAttr := s.Children().Eq(1).Find("a").Attr("href"); hasAttr {
				val = strings.TrimLeft(val, "/user/")
				val = strings.TrimRight(val, "/")
				result.UploadedBy = val
			}
		case "Downloads":
			result.Downloads, _ = strconv.Atoi(value)
		case "Date uploaded":
			result.DateUploaded = value
		case "Seeders":
			result.Seeders, _ = strconv.Atoi(value)
		case "Leechers":
			result.Leechers, _ = strconv.Atoi(value)
		case "Last checked":
			result.LastChecked = value
		}
	})

	return
}

func (doc *document) getDetailsPageImage() (image *url.URL, err error) {
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

func (doc *document) getDetailsPageFilm() (
	film struct {
		title string
		url   *url.URL
	},
	err error,
) {
	const selector = ".torrent-detail-info h3 a"
	selection := doc.Find(selector)
	if selection.Nodes == nil {
		err = fmt.Errorf("missing film info (title/link) at selector '%s'", selector)
		return
	}
	film.title = selection.Text()

	href, exists := selection.Attr("href")
	if !exists {
		return film, errors.New("film title link has no href")
	}

	film.url, err = doc.Url.Parse(href)

	return
}

func (doc *document) getDetailsPageFilmGenres() (categories []string, err error) {
	const selector = ".torrent-category span"
	selection := doc.Find(selector)
	if selection.Nodes == nil {
		err = fmt.Errorf("missing film categories at selector '%s'", selector)
		return
	}
	selection.Each(func(_ int, s *goquery.Selection) {
		categories = append(categories, s.Text())
	})
	return
}

func (doc *document) getDetailsPageFilmDescription() (description string, err error) {
	const selector = ".torrent-detail-info p"
	selection := doc.Find(selector)
	if selection.Nodes == nil {
		err = fmt.Errorf("missing film description element at selector '%s'", selector)
		return
	}
	description = strings.TrimSpace(selection.Text())
	return
}

func (doc *document) getDetailsPageImdb(re *regexp.Regexp) (imdb *url.URL, err error) {
	html, _ := doc.Html()
	if matches := re.FindAllString(html, -1); matches != nil {
		imdb, _ = url.Parse(matches[0])
		return
	}
	err = errors.New("couldn't find an imdb link in page")
	return
}

func (doc *document) GetData() *Torrent {
	t := new(Torrent)
	t.FoundOn = doc.Url

	var err error

	if t.ID, err = doc.getID(); err != nil {
		log.WithError(err).Warn("missing torrent id")
	}

	if t.Title, err = doc.getDetailsPageTitle(); err != nil {
		log.WithError(err).Fatal("missing title element")
	}

	t.TitleInfo = coll33tx.ParseTitleInfo(t.Title)

	t.Quality = t.TitleInfo.Quality

	if t.Year = doc.getYear(); t.Year == 0 {
		t.Year = t.TitleInfo.Year
	}

	if t.Magnet, err = doc.getDetailsPageMagnet(); err != nil {
		log.WithError(err).Debug("missing magnet")
	}

	if box, err := doc.getDetailsPageBox(); err == nil {
		t.Category = box.Category
		t.Type = box.Type
		t.Language = box.Language
		t.TotalSize = box.TotalSize
		t.UploadedBy = box.UploadedBy
		t.Downloads = box.Downloads
		t.LastChecked = box.LastChecked
		t.DateUploaded = box.DateUploaded
		t.Seeders = box.Seeders
		t.Leechers = box.Leechers
	} else {
		log.WithError(err).Warn()
	}

	if t.Image, err = doc.getDetailsPageImage(); err != nil {
		log.WithError(err).Debug()
	}

	film, err := doc.getDetailsPageFilm()
	if err != nil {
		log.WithError(err).Debug("missing film info box")
		t.FilmCleanTitle = t.TitleInfo.Title
	} else {
		t.FilmCleanTitle = film.title
		t.FilmLink = film.url
	}

	if t.Genres, err = doc.getDetailsPageFilmGenres(); err != nil {
		log.WithError(err).Debug()
	}

	if t.FilmDescription, err = doc.getDetailsPageFilmDescription(); err != nil {
		log.WithError(err).Debug()
	}

	t.IMDB, _ = doc.getDetailsPageImdb(imdbRE)

	return t
}

// todo test
func (doc *document) getYear() int {
	title, err := doc.getDetailsPageTitle()
	if err != nil {
		return 0
	}

	year, err := strconv.Atoi(yearRE.FindString(title))
	if err != nil {
		return 0
	}

	return year
}

func (doc *document) getID() (int, error) {
	matches := idRE.FindStringSubmatch(doc.Url.Path)

	id, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, errors.New("Coult not get torrent ID")
	}

	return id, nil
}
