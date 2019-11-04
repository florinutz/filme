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
		goIntoDetails        bool
		debugLevel           value.DebugLevelValue
	}

	cmd = &cobra.Command{
		Use:   "search <what>",
		Short: "Search torrents",

		RunE: func(cmd *cobra.Command, args []string) error {
			return f.Search(
				args,
				opts.numberOfDesiredItems,
				opts.goIntoDetails,
				opts.debugLevel,
			)
		},
	}

	cmd.Flags().IntVarP(&opts.numberOfDesiredItems, "wanted-items", "n", 20,
		"keep fetching pages until the number of result items was met")
	cmd.Flags().BoolVarP(&opts.goIntoDetails, "crawl-details", "d", false,
		"follows every link in the list and fetches detail pages data")
	defaultDebugLevel := logrus.DebugLevel
	_ = opts.debugLevel.Set(defaultDebugLevel.String())
	cmd.Flags().Var(&opts.debugLevel, "debug-level", fmt.Sprintf("possible debug levels: %s",
		strings.Join(value.GetAllLevels(), ", ")))

	return
}
