package filter

import "github.com/spf13/pflag"

type Filter struct {
	MaxItems uint
	Seeders  *uintVal
	Leechers *uintVal
	Size     *uintVal
}

type uintVal struct {
	Min uint
	Max uint
}

func (f *Filter) GetLinkedFlagSet() *pflag.FlagSet {
	set := pflag.NewFlagSet("filters", pflag.ExitOnError)

	if f.Seeders == nil {
		f.Seeders = new(uintVal)
	}
	if f.Leechers == nil {
		f.Leechers = new(uintVal)
	}
	if f.Size == nil {
		f.Size = new(uintVal)
	}

	set.UintVarP(&f.MaxItems, "total", "t", 50, "specifies the maximum desired number of items to display")

	set.UintVar(&f.Seeders.Max, "max-seeders", 0, "only show values below this")
	set.UintVar(&f.Seeders.Min, "min-seeders", 0, "only show values above this")

	set.UintVar(&f.Leechers.Max, "max-leechers", 0, "only show values below this")
	set.UintVar(&f.Leechers.Min, "min-leechers", 0, "only show values above this")

	set.UintVar(&f.Size.Max, "max-size", 0, "only show values below this")

	return set
}
