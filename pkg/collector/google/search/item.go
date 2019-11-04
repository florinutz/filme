package search

import (
	"fmt"
	"net/url"
	"strconv"
)

type (
	// BaseItem represents a google search result item
	BaseItem interface {
		Title() string
		SetTitle(string)

		Url() *url.URL
		SetUrl(*url.URL)

		Description() string
		SetDescription(string)

		Offset() int
		SetOffset(int)

		Errors() []error
		SetErrors([]error)
		AddErrors([]error)
	}

	WithRating interface {
		BaseItem

		Acts() int
		SetActs(acts int)

		ActType() string
		SetActType(actType string)

		Rating() int
		SetRating(rating int)

		RatingStr() string
	}

	// ItemDefault implements BaseItem
	ItemDefault struct {
		title       string
		url         *url.URL
		description string
		offset      int
		errs        []error
	}

	// ItemWithRating adds rating to ItemDefault (int out of 100)
	ItemWithRating struct {
		BaseItem
		rating  int
		acts    int
		actType string
	}

	ItemNetflix struct {
		*ItemDefault
	}

	ItemPrime struct {
		*ItemDefault
	}

	ItemWikipedia struct {
		*ItemDefault
	}

	ItemImdb struct {
		*ItemWithRating
	}

	ItemMetacritic struct {
		*ItemWithRating
	}

	ItemRottenTomatoes struct {
		*ItemWithRating
	}
)

// Decorate checks the url and returns either a decorated item (Imdb, Netflix, etc) or nil
func (i *ItemDefault) Decorate() (result BaseItem) {
	if i.Url() == nil {
		return
	}

	switch i.Url().Hostname() {
	case "www.netflix.com":
		result = &ItemNetflix{i}
	case "www.amazon.com":
		result = &ItemPrime{i}
	case "www.imdb.com":
		result = &ItemImdb{&ItemWithRating{BaseItem: i}}
	case "en.wikipedia.org":
		result = &ItemWikipedia{i}
	case "www.metacritic.com":
		result = &ItemMetacritic{&ItemWithRating{BaseItem: i}}
	case "www.rottentomatoes.com":
		result = &ItemRottenTomatoes{&ItemWithRating{BaseItem: i}}
	}

	return
}

func (iwr *ItemWithRating) ActType() string {
	return iwr.actType
}

func (iwr *ItemWithRating) SetActType(actType string) {
	iwr.actType = actType
}

func (iwr *ItemWithRating) Acts() int {
	return iwr.acts
}

func (iwr *ItemWithRating) SetActs(acts int) {
	iwr.acts = acts
}

func (iwr *ItemWithRating) Rating() int {
	return iwr.rating
}

func (iwr *ItemWithRating) SetRating(rating int) {
	iwr.rating = rating
}

func (i *ItemDefault) SetErrors(errs []error) {
	i.errs = errs
}

func (i *ItemDefault) SetOffset(offset int) {
	i.offset = offset
}

func (i *ItemDefault) SetDescription(description string) {
	i.description = description
}

func (i *ItemDefault) SetUrl(url *url.URL) {
	i.url = url
}

func (i *ItemDefault) SetTitle(title string) {
	i.title = title
}

func (i *ItemDefault) Title() string {
	return i.title
}

func (i *ItemDefault) Url() *url.URL {
	return i.url
}

func (i *ItemDefault) Description() string {
	return i.description
}

func (i *ItemDefault) Offset() int {
	return i.offset
}

func (i *ItemDefault) Errors() []error {
	return i.errs
}

func (i *ItemDefault) AddErrors(errs []error) {
	i.errs = append(i.errs, errs...)
}

func (iwr *ItemWithRating) RatingStr() string {
	if iwr.Rating() == 0 {
		return ""
	}
	return fmt.Sprintf("%d%% from %d %s", iwr.Rating(), iwr.Acts(), iwr.ActType())
}

func (i ItemImdb) RatingStr() string {
	if i.Rating() == 0 {
		return ""
	}
	rating := strconv.FormatFloat(float64(i.Rating())/float64(10), 'f', 1, 32)
	return fmt.Sprintf("%s from %d %s", rating, i.Acts(), i.ActType())
}
