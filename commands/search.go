package commands

import (
	"fmt"
	"strings"

	"github.com/florinutz/filme/pkg/config/value"
	"github.com/florinutz/filme/pkg/config/value/1337x/list/encoding"
	"github.com/florinutz/filme/pkg/config/value/1337x/list/search_category"
	"github.com/florinutz/filme/pkg/config/value/1337x/list/sort"
	"github.com/florinutz/filme/pkg/filme"
	"github.com/florinutz/filme/pkg/filme/l33tx/list/filter"
	"github.com/florinutz/filme/pkg/filme/l33tx/list/input"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// BuildSearchCmd mirrors the 1337x search cmd
func BuildSearchCmd(f *filme.Filme) (cmd *cobra.Command) {
	var opts struct {
		goIntoDetails bool // todo implement this
		debugLevel    value.DebugLevelValue
		sort          sort.Value
		category      search_category.SearchCategory
		filters       filter.Filter
		encoding      encoding.ListEncoding
	}

	cmd = &cobra.Command{
		Use:   "search <what>",
		Short: "Search torrents",

		RunE: func(cmd *cobra.Command, args []string) error {
			inputs := input.ListingInput{
				Search:   strings.Join(args, " "),
				Url:      nil, // todo merge search and the other command
				Category: &opts.category,
				Encoding: &opts.encoding,
				Sort:     opts.sort,
			}

			return f.Search(opts.goIntoDetails, inputs, opts.filters, opts.debugLevel)
		},
	}

	filters := &opts.filters
	cmd.Flags().AddFlagSet(filters.GetLinkedFlagSet())

	cmd.Flags().BoolVarP(&opts.goIntoDetails, "crawl-details", "d", false,
		"follows every link in the list and fetches detail pages data")

	// default debug level
	defaultDebugLevel := logrus.DebugLevel
	_ = opts.debugLevel.Set(defaultDebugLevel.String())
	cmd.Flags().VarP(&opts.debugLevel, "debug-level", "l", fmt.Sprintf("one of: %s",
		strings.Join(value.GetAllLevels(), ", ")))

	// default movie category "all"
	opts.category = search_category.SearchCategoryAll
	cmd.Flags().VarP(&opts.category, "category", "c", fmt.Sprintf("one of: %s",
		strings.Join(search_category.GetAllSearchCategories(), ", ")))

	// default movie encoding 1080p
	opts.encoding = encoding.EncHD
	cmd.Flags().VarP(&opts.encoding, "encoding", "e", fmt.Sprintf("one of: %s",
		strings.Join(encoding.GetAll(), ", ")))

	// default sorting
	defaultSort, err := sort.NewValue("time-desc")
	if err != nil {
		panic("shouldn't happen, since the value above is valid. RIGHT?")
	}
	opts.sort = *defaultSort
	cmd.Flags().VarP(&opts.sort, "sort", "s", fmt.Sprintf("one of: %s",
		strings.Join(sort.GetAllValues(), ", ")))

	return
}
