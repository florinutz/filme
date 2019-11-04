package commands

import (
	"fmt"
	"strings"

	"github.com/florinutz/filme/pkg/config"

	"github.com/sirupsen/logrus"

	"github.com/florinutz/filme/pkg/config/value"
	"github.com/florinutz/filme/pkg/filme"
	"github.com/spf13/cobra"
)

func BuildGoogleCmd(f *filme.Filme) *cobra.Command {
	var opts struct {
		debugLevel           value.DebugLevelValue
		url                  string
		resultsCount         int
		language             string
		onlyFilmRelatedItems bool
	}

	cmd := &cobra.Command{
		Use:   "google <search|url>",
		Short: "Searches google",

		PreRun: func(cmd *cobra.Command, args []string) {
			opts.url, _ = config.GetGoogleUrlFromArgs(args, opts.resultsCount, opts.language)
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			return f.SearchGoogle(opts.url, opts.onlyFilmRelatedItems, opts.debugLevel)
		},
	}

	defaultDebugLevel := logrus.DebugLevel
	_ = opts.debugLevel.Set(defaultDebugLevel.String())
	cmd.Flags().Var(&opts.debugLevel, "debug-level", fmt.Sprintf("possible debug levels: %s",
		strings.Join(value.GetAllLevels(), ", ")))

	cmd.Flags().IntVarP(&opts.resultsCount, "results-number", "n", 10,
		"the desired number of results")

	cmd.Flags().StringVarP(&opts.language, "language", "l", "lang_en",
		"the desired language of the results")

	cmd.Flags().BoolVarP(&opts.onlyFilmRelatedItems, "only-film-related", "f", true,
		"only show results related to films (metacritic, imdb, rottentomatoes, wikipedia, etc)")

	return cmd
}
