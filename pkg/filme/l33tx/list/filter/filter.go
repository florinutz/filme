package filter

import "github.com/spf13/pflag"

type Filter struct {
	MaxItems uint
	Seeders  *UintVal
	Leechers *UintVal
	Pages    *UintVal
	Size     *UintVal
}

type UintVal struct {
	Min uint
	Max uint
}

func (f *Filter) GetLinkedFlagSet() *pflag.FlagSet {
	set := pflag.NewFlagSet("filters", pflag.ExitOnError)

	if f.Seeders == nil {
		f.Seeders = new(UintVal)
	}
	if f.Leechers == nil {
		f.Leechers = new(UintVal)
	}
	if f.Pages == nil {
		f.Pages = new(UintVal)
	}
	if f.Size == nil {
		f.Size = new(UintVal)
	}

	set.UintVarP(&f.MaxItems, "total", "t", 20, "specifies the maximum desired number of "+
		"items to display.\nDefaults to one page's worth of items.")

	set.UintVar(&f.Seeders.Max, "seeders-max", 0, "ignores items with more seeders")
	set.UintVar(&f.Seeders.Min, "seeders-min", 0, "ignores items with less seeders")

	set.UintVar(&f.Leechers.Max, "leechers-max", 0, "ignores items with more leechers")
	set.UintVar(&f.Leechers.Min, "leechers-min", 0, "ignores items with less leechers")

	set.UintVar(&f.Pages.Max, "page-max", 0, "stop at this page")
	set.UintVar(&f.Pages.Min, "page-min", 0, "start at this page")

	set.UintVar(&f.Size.Max, "size", 0, "ignore items bigger than this")

	return set
}

func (f *Filter) GetLogFields() (result map[string]interface{}) {
	result = map[string]interface{}{}
	result["filter_max_items"] = f.MaxItems
	if f.Seeders != nil {
		if f.Seeders.Min != 0 {
			result["filter_seeders_min"] = f.Seeders.Min
		}
		if f.Seeders.Max != 0 {
			result["filter_seeders_max"] = f.Seeders.Max
		}
	}
	if f.Leechers != nil {
		if f.Leechers.Min != 0 {
			result["filter_leechers_min"] = f.Leechers.Min
		}
		if f.Leechers.Max != 0 {
			result["filter_leechers_max"] = f.Leechers.Max
		}
	}
	if f.Pages != nil {
		if f.Pages.Min != 0 {
			result["filter_pages_min"] = f.Pages.Min
		}
		if f.Pages.Max != 0 {
			result["filter_pages_max"] = f.Pages.Max
		}
	}
	if f.Size != nil {
		if f.Size.Min != 0 {
			result["filter_size_min"] = f.Size.Min
		}
		if f.Size.Max != 0 {
			result["filter_size_max"] = f.Size.Max
		}
	}
	return
}
