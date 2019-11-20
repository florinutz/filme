package commands

import (
	"fmt"
	"strings"

	"github.com/florinutz/filme/pkg/config/value"
	"github.com/florinutz/filme/pkg/filme"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// BuildSearchCmd mirrors the 1337x search cmd
func BuildSearchCmd(f *filme.Filme) (cmd *cobra.Command) {
	var opts struct {
		numberOfDesiredItems int
		goIntoDetails        bool // todo implement this
		debugLevel           value.DebugLevelValue
		sort                 value.LeetxListSortValue
		category             value.LeetxListSearchCategory
		encoding             value.LeetxListEncoding
	}

	cmd = &cobra.Command{
		Use:   "search <what>",
		Short: "Search torrents",

		RunE: func(cmd *cobra.Command, args []string) error {
			return f.Search(
				strings.Join(args, " "),
				opts.numberOfDesiredItems,
				opts.goIntoDetails,
				opts.category,
				opts.encoding,
				opts.sort,
				opts.debugLevel,
			)
		},
	}

	cmd.Flags().IntVarP(&opts.numberOfDesiredItems, "wanted-items", "n", 20,
		"keep fetching pages until the number of result items was met")
	cmd.Flags().BoolVarP(&opts.goIntoDetails, "crawl-details", "d", false,
		"follows every link in the list and fetches detail pages data")

	// default debug level
	defaultDebugLevel := logrus.DebugLevel
	_ = opts.debugLevel.Set(defaultDebugLevel.String())
	cmd.Flags().VarP(&opts.debugLevel, "debug-level", "l", fmt.Sprintf("one of: %s",
		strings.Join(value.GetAllLevels(), ", ")))

	// default movie category "all"
	opts.category = value.SearchCategoryAll
	cmd.Flags().VarP(&opts.category, "category", "c", fmt.Sprintf("one of: %s",
		strings.Join(value.GetAllLeetxListMovieSearchCategoryValues(), ", ")))

	// default movie encoding 1080p
	opts.encoding = value.EncodingHD
	cmd.Flags().VarP(&opts.encoding, "encoding", "e", fmt.Sprintf("one of: %s",
		strings.Join(value.GetAllLeetxListMovieEncodingValues(), ", ")))

	// default sorting
	defaultSort, err := value.NewLeetxListSortValue("time-desc")
	if err != nil {
		panic("shouldn't happen, since the value above is valid. RIGHT?")
	}
	opts.sort = *defaultSort
	cmd.Flags().VarP(&opts.sort, "sort", "s", fmt.Sprintf("one of: %s",
		strings.Join(value.GetAllLeetxListSortValues(), ", ")))

	return
}
