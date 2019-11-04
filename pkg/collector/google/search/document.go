package search

import (
	"fmt"
	"math"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/florinutz/filme/pkg/collector"
	"github.com/pkg/errors"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

const (
	TestSearch = "https://www.google.com/search?q=imdb|rottentomatoes|metacritic|netflix|wikipedia+love+death+robots&num=30"

	groupRating   = "rating"
	groupActs     = "acts"
	groupActsType = "act_type"
)

var ratingRegex *regexp.Regexp

func init() {
	ratingRegex = regexp.MustCompile(fmt.Sprintf(
		`(?P<%s>\d+(([,.]?\d*/\d+)|%%))\D+(?P<%s>\d+[,.]?\d*)\W+(?P<%s>.+)?`,
		groupRating, groupActs, groupActsType,
	))

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

func (doc *document) GetItems() (map[int]BaseItem, error) {
	selection := doc.Find(".g")
	if selection.Nodes == nil {
		return nil, errors.New("couldn't find the google results container")
	}

	items := make(map[int]BaseItem)

	selection.Each(doc.itemExtractor(items))

	return items, nil
}

func (doc *document) itemExtractor(items map[int]BaseItem) func(i int, listItem *goquery.Selection) {
	return func(i int, listItem *goquery.Selection) {
		var item BaseItem

		item = &ItemDefault{offset: i}

		desc := listItem.Find(".s .st")
		if desc.Nodes == nil {
			item.SetErrors(append(item.Errors(), errors.New("missing description")))
		}
		item.SetDescription(desc.Text())

		titleH3 := listItem.Find("h3")
		if titleH3.Nodes == nil {
			err := errors.New("missing title")
			item.AddErrors([]error{err})
		}
		item.SetTitle(titleH3.Text())

		resultLink := listItem.Find(".r>a")
		if resultLink.Nodes == nil {
			err := errors.New("missing result link")
			item.AddErrors([]error{err})
		}

		outboundUrlStr, exists := resultLink.Attr("href")
		if !exists {
			err := errors.New("missing href")
			item.AddErrors([]error{err})
		}

		outboundUrl, err := url.Parse(outboundUrlStr)
		if err != nil {
			err = errors.Wrap(err, "could not parse google outbound relative url")
			item.AddErrors([]error{err})
		}

		if outboundUrl.IsAbs() {
			item.SetUrl(outboundUrl)
		} else {
			actualUrl, err := url.Parse(outboundUrl.Query().Get("q"))
			if err != nil {
				err = errors.Wrap(err, "could not parse the q parameter outbound url")
				item.AddErrors([]error{err})
			} else {
				item.SetUrl(actualUrl)
			}
		}

		if len(item.Errors()) > 2 {
			return
		}

		itemDefault, _ := item.(*ItemDefault)
		if decorated := itemDefault.Decorate(); decorated != nil {
			item = decorated
		}

		if rating := getRating(listItem); rating != nil {
			if ratedItem, ok := item.(WithRating); ok {
				ratedItem.SetActs(rating.acts)
				ratedItem.SetActType(rating.actType)
				ratedItem.SetRating(rating.rating)
				item = ratedItem
			} else {
				doc.log.Warn("rating found, but the current decoration doesn't support it")
			}
		}

		items[i] = item
	}
}

// todo write tests for this
// getRating returns an enriched (WithRating) version of the previous BaseItem
func getRating(listItem *goquery.Selection) *struct {
	rating, acts int
	actType      string
	errs         []error
} {
	starsTag := listItem.Find("g-review-stars")
	if starsTag.Nodes == nil {
		return nil
	}
	ratingText := starsTag.Parent().Text() // Rating: 8,8/10 - ‎39,505 votes
	groupNames := ratingRegex.SubexpNames()
	matches := ratingRegex.FindAllStringSubmatch(ratingText, -1)
	if len(matches) == 0 {
		return nil
	}

	result := new(struct {
		rating, acts int
		actType      string
		errs         []error
	})

	for _, match := range matches {
		for groupIdx, value := range match {
			switch groupNames[groupIdx] {
			case groupRating: //  comes as 7,3/10 (imdb) or 73% (rt)
				var (
					intVal int
					err    error
				)
				if strings.Contains(value, "/") {
					fraction := strings.Split(value, "/")
					if len(fraction) != 2 {
						err = fmt.Errorf("value '%s' is not a fraction", value)
						result.errs = append(result.errs, err)
						continue
					}
					if fraction[1] != "10" {
						err = fmt.Errorf("expecting /10, got /%s", fraction[1])
						result.errs = append(result.errs, err)
						continue
					}
					value := fraction[0] // 7,3
					value = strings.Replace(value, ",", ".", 1)
					f, err := strconv.ParseFloat(value, 32)
					if err != nil {
						err := errors.Wrapf(err, "can't parse float '%s'", fraction[1])
						result.errs = append(result.errs, err)
						continue
					}
					intVal = int(math.Round(f * 10)) // only keep 1 decimal
				} else if value[len(value)-1:] == "%" {
					intVal, err = strconv.Atoi(strings.TrimRight(value, "%"))
					if err != nil {
						err := errors.Wrapf(err, "can't get percentage from '%s'", value)
						result.errs = append(result.errs, err)
						continue
					}
				} else {
					err := errors.Wrapf(err, "weird rating string '%s'", value)
					result.errs = append(result.errs, err)
					continue
				}
				result.rating = intVal
			case groupActs:
				value = strings.Replace(value, ",", "", 1)
				intVal, err := strconv.Atoi(value)
				if err != nil {
					err = errors.Wrap(err, "can't get rating's number of interactions")
					result.errs = append(result.errs, err)
					continue
				}
				result.acts = intVal
			case groupActsType:
				value = strings.TrimLeft(value, " ")
				if value == "" {
					err := errors.New("can't get rating's interactions type")
					result.errs = append(result.errs, err)
					continue
				}
				result.actType = value
			}
		}
	}

	return result
}
