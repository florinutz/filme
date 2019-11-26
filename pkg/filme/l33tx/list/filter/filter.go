package filter

import "github.com/spf13/pflag"

type Filter struct {
	MaxItems uint
	Seeders  *UintVal
	Leechers *UintVal
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
	if f.Size == nil {
		f.Size = new(UintVal)
	}

	set.UintVarP(&f.MaxItems, "total", "t", 50, "specifies the maximum desired number of items to display")

	set.UintVar(&f.Seeders.Max, "max-seeders", 0, "only show values below this")
	set.UintVar(&f.Seeders.Min, "min-seeders", 0, "only show values above this")

	set.UintVar(&f.Leechers.Max, "max-leechers", 0, "only show values below this")
	set.UintVar(&f.Leechers.Min, "min-leechers", 0, "only show values above this")

	set.UintVar(&f.Size.Max, "max-size", 0, "only show values below this")

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
