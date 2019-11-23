package list

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	"github.com/florinutz/filme/pkg/collector"
	"github.com/florinutz/filme/pkg/filme/l33tx/list/filter"
	"github.com/gocolly/colly"
)

const (
	TestPageList               = "https://1337x.to/popular-movies"
	TestPageListWithPagination = "https://1337x.to/search/romania/3/"
)

// document describes a list page
type document struct {
	*goquery.Document
}

func NewDocument(r *colly.Response) (*document, error) {
	d, err := collector.GetResponseDocument(r)
	if err != nil {
		return nil, err
	}
	return &document{d}, nil
}

type Item struct {
	Name     string
	Href     *url.URL
	Size     string
	Seeders  int
	Leechers int
}

// Validate returns all the problems related to a list item
func (item *Item) Validate(f filter.Filter) (errs []error) {
	if f.Size != nil {
		size, err := ToBytes(item.Size)
		if err != nil {
			if f.Size.Max > 0 && size > uint64(f.Size.Max) {
				errs = append(errs, fmt.Errorf("item size %d bigger than max allowed size %d", size, f.Size.Max))
			}
			if f.Size.Min > 0 && size < uint64(f.Size.Min) {
				errs = append(errs, fmt.Errorf("item size %d smaller than max allowed size %d", size, f.Size.Min))
			}
		}
	}
	if f.Seeders.Max > 0 && int(f.Seeders.Max) < item.Seeders {
		errs = append(errs, fmt.Errorf("seeders count %d bigger than the allowed %d max", item.Seeders, f.Seeders.Max))
	}
	if f.Seeders.Min > 0 && int(f.Seeders.Min) > item.Seeders {
		errs = append(errs, fmt.Errorf("seeders count %d smaller than the allowed %d min", item.Seeders, f.Seeders.Min))
	}

	if f.Leechers.Max > 0 && int(f.Leechers.Max) < item.Leechers {
		errs = append(errs, fmt.Errorf("leechers count %d bigger than the allowed %d max", item.Leechers, f.Leechers.Min))
	}
	if f.Leechers.Min > 0 && int(f.Leechers.Min) > item.Leechers {
		errs = append(errs, fmt.Errorf("leechers count %d smaller than the allowed %d min", item.Leechers, f.Leechers.Min))
	}

	return
}

// ToBytes parses a string formatted by ByteSize as bytes
func ToBytes(s string) (uint64, error) {
	err := errors.New("invalid")

	s = strings.TrimSpace(s)
	s = strings.ToUpper(s)

	i := strings.IndexFunc(s, unicode.IsLetter)
	if i == -1 {
		return 0, err
	}

	bytesString, multiple := s[:i], s[i:]
	bytesString = strings.Trim(bytesString, " ")
	bytesString = strings.Replace(bytesString, ",", ".", 1)
	bytes, err := strconv.ParseFloat(bytesString, 64)
	if err != nil || bytes < 0 {
		return 0, err
	}

	const (
		BYTE = 1 << (10 * iota)
		KILOBYTE
		MEGABYTE
		GIGABYTE
		TERABYTE
		PETABYTE
		EXABYTE
	)

	switch multiple {
	case "E", "EB", "EIB":
		return uint64(bytes * EXABYTE), nil
	case "P", "PB", "PIB":
		return uint64(bytes * PETABYTE), nil
	case "T", "TB", "TIB":
		return uint64(bytes * TERABYTE), nil
	case "G", "GB", "GIB":
		return uint64(bytes * GIGABYTE), nil
	case "M", "MB", "MIB":
		return uint64(bytes * MEGABYTE), nil
	case "K", "KB", "KIB":
		return uint64(bytes * KILOBYTE), nil
	case "B":
		return uint64(bytes), nil
	default:
		return 0, err
	}
}

type Line struct {
	Item *Item
	Errs []error
}

// GetLines returns list items along with their errors / missing stuff
func (doc *document) GetLines() ([]*Line, error) {
	// look for no results msg
	possibleNotFound := doc.Find(".page-content .box-info-detail p").Text()
	if strings.Trim(possibleNotFound, " ") == "No results were returned." {
		return nil, fmt.Errorf("no results")
	}

	// look for "smth went wrong" title
	possibleErrorTitle := doc.Find("title").Text()
	if strings.Trim(possibleErrorTitle, " ") == "Error something went wrong." {
		return nil, fmt.Errorf("html title says smth went wrong")
	}

	// look for the actual lines
	selector := ".page-content .table-striped tbody tr"
	trs := doc.Find(selector)
	if trs.Nodes == nil {
		return nil, fmt.Errorf("selector '%s' not found in the retrieved html", selector)
	}

	var lines []*Line
	trs.Each(func(i int, tr *goquery.Selection) {
		line := new(Line)
		line.Item, line.Errs = doc.trToItem(i, tr)
		lines = append(lines, line)
	})

	return lines, nil
}

func (doc *document) trToItem(i int, tr *goquery.Selection) (item *Item, errs []error) {
	item = new(Item)
	var err error

	td := tr.Find("td.name")
	if td.Nodes != nil {
		item.Name = td.Text()
	}
	aName := td.Find("a").Eq(1)
	if aName.Nodes != nil {
		href, exists := aName.Attr("href")
		if !exists {
			errs = append(errs, fmt.Errorf("list item %d (%s) has no link behind", i, item.Name))
		}
		item.Href, err = doc.Url.Parse(href)
		if err != nil {
			errs = append(errs, fmt.Errorf("link behind link %d (%s) is invalid", i, item.Name))
		}
	}

	td = tr.Find("td.seeds")
	if td.Nodes != nil {
		if item.Seeders, err = strconv.Atoi(td.Text()); err != nil {
			errs = append(errs, fmt.Errorf("can't convert seeders to int for item %d (%s)", i, item.Name))
		}
	}

	td = tr.Find("td.leeches")
	if td.Nodes != nil {
		if item.Leechers, err = strconv.Atoi(td.Text()); err != nil {
			errs = append(errs, fmt.Errorf("can't convert leechers to int for item %d (%s)", i, item.Name))
		}
	}

	td = tr.Find("td.size")
	if td.Nodes != nil {
		unwantedSpan := td.Find("span")
		if unwantedSpan.Nodes != nil {
			unwantedSpan.Remove()
		}
		item.Size = td.Text()
		if utf8.RuneCountInString(item.Size) < 3 {
			errs = append(errs, fmt.Errorf("weird torrent size for item %d (%s)", i, item.Name))
		}
	}

	return item, errs
}
