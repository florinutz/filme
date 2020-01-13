package commands

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/florinutz/filme/pkg/config/value/1337x/list/encoding"
	"github.com/florinutz/filme/pkg/config/value/1337x/list/search_category"
	"github.com/florinutz/filme/pkg/config/value/1337x/list/sort"
	"github.com/florinutz/filme/pkg/filme"
	"github.com/florinutz/filme/pkg/filme/l33tx/list/filter"
	"github.com/florinutz/filme/pkg/filme/l33tx/list/input"
	"github.com/spf13/cobra"
)

// BuildSearchCmd mirrors the 1337x search cmd
func BuildSearchCmd(f *filme.Filme) (cmd *cobra.Command) {
	var opts struct {
		goIntoDetails bool // todo implement this
		sort          sort.Value
		category      search_category.SearchCategory
		filters       filter.Filter
		encoding      encoding.ListEncoding
	}

	cmd = &cobra.Command{
		Use:   "search <what>",
		Short: "Search torrents",

		RunE: func(cmd *cobra.Command, args []string) error {
			what := strings.Join(args, " ")
			inputs := input.ListingInput{
				Category: &opts.category,
				Encoding: &opts.encoding,
				Sort:     opts.sort,
			}

			if strings.HasPrefix(what, "https://1337x.to/") {
				var err error
				if inputs.URL, err = url.Parse(what); err != nil {
					return err
				}
			} else {
				inputs.Search = what
			}

			return f.Search(opts.goIntoDetails, inputs, opts.filters)
		},
	}

	filters := &opts.filters
	cmd.Flags().AddFlagSet(filters.GetLinkedFlagSet())

	cmd.Flags().BoolVarP(&opts.goIntoDetails, "crawl-details", "d", false,
		"follows every link in the list and fetches detail pages data")

	// default movie category "all"
	opts.category = search_category.SearchCategoryAll
	cmd.Flags().VarP(&opts.category, "category", "c", fmt.Sprintf("one of: %s",
		strings.Join(search_category.GetAllSearchCategories(), ", ")))

	cmd.Flags().VarP(&opts.encoding, "encoding", "e", fmt.Sprintf("one of: %s",
		strings.Join(encoding.GetAll(), ", ")))

	// default sorting
	defaultSort, err := sort.NewValue("seeders-desc")
	if err != nil {
		panic("shouldn't happen, since the value above is valid. RIGHT?")
	}
	opts.sort = *defaultSort
	cmd.Flags().VarP(&opts.sort, "sort", "s", fmt.Sprintf("one of: %s",
		strings.Join(sort.GetAllValues(), ", ")))

	return
}
