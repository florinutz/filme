package line

import (
	"fmt"
	"net/url"

	"github.com/florinutz/filme/pkg/filme/l33tx/list/filter"
)

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
