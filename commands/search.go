package commands

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/florinutz/filme/pkg/config/value/1337x/list/encoding"
	"github.com/florinutz/filme/pkg/config/value/1337x/list/search_category"
	"github.com/florinutz/filme/pkg/config/value/1337x/list/sort"
	"github.com/florinutz/filme/pkg/filme"
	"github.com/florinutz/filme/pkg/filme/l33tx/list"
	"github.com/florinutz/filme/pkg/filme/l33tx/list/filter"
	"github.com/florinutz/filme/pkg/filme/l33tx/list/input"
	"github.com/spf13/cobra"
)

// BuildSearchCmd mirrors the 1337x search cmd
func BuildSearchCmd(f *filme.Filme) (cmd *cobra.Command) {
	var opts struct {
		goIntoDetails                   bool // todo implement this
		sort                            sort.Value
		category                        search_category.SearchCategory
		filters                         filter.Filter
		encoding                        encoding.ListEncoding
		delay, randomDelay, parallelism int
		userAgent                       string
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

			logFields := getSearchLogFields(inputs, opts.goIntoDetails, opts.filters)
			log := *f.Log.WithFields(logFields)

			ls := list.NewList(inputs, opts.filters, log)

			f.Log = log.Logger
			// populates the list
			if err := f.Search(ls, opts.goIntoDetails, opts.delay, opts.randomDelay, opts.parallelism, opts.userAgent); err != nil {
				return err
			}

			ls.Display(f.Out)

			return nil
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

	cmd.Flags().IntVar(&opts.delay, "reqs-delay", 3, "requests delay")
	cmd.Flags().IntVar(&opts.randomDelay, "reqs-random-delay", 3, "requests random delay")
	cmd.Flags().IntVar(&opts.parallelism, "reqs-parallelism", 4, "requests parallelism")
	cmd.Flags().StringVar(&opts.userAgent, "user-agent", "", "requests user agent. Leave this empty for random browser UAs")

	return
}

func getSearchLogFields(inputs input.ListingInput, goIntoDetails bool, filters filter.Filter) (result map[string]interface{}) {
	result = map[string]interface{}{
		"search_inputs_sort": inputs.Sort,
		"with_details":       goIntoDetails,
	}
	if inputs.Encoding != nil && *inputs.Encoding != 0 {
		result["search_inputs_encoding"] = *inputs.Encoding
	}
	if inputs.Category != nil && *inputs.Category != 0 {
		result["search_inputs_category"] = *inputs.Category
	}
	if inputs.Search != "" {
		result["search"] = inputs.Search
	}
	if inputs.URL != nil {
		result["search_inputs_url"] = inputs.URL.String()
	}

	for key, v := range filters.GetLogFields() {
		result[key] = v
	}

	return
}
