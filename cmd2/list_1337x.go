package main

import (
	"fmt"
	"strings"

	"github.com/florinutz/filme/pkg/config"

	"github.com/sirupsen/logrus"

	"github.com/florinutz/filme/pkg/config/value"
	"github.com/florinutz/filme/pkg/filme"
	"github.com/spf13/cobra"
)

func build1337xParseListPageCmd(f *filme.Filme) *cobra.Command {
	var opts struct {
		goIntoDetails bool
		debugLevel    value.DebugLevelValue
		url           string
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "Parses 1337x listings",

		PreRun: func(cmd *cobra.Command, args []string) {
			opts.url = config.GetListUrlFromArgs(args)
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			return f.VisitList(opts.url, opts.goIntoDetails, opts.debugLevel)
		},
	}

	cmd.Flags().BoolVarP(&opts.goIntoDetails, "crawl-details", "d", false,
		"follows every link in the list and fetches detail pages data")
	defaultDebugLevel := logrus.DebugLevel
	_ = opts.debugLevel.Set(defaultDebugLevel.String())
	cmd.Flags().Var(&opts.debugLevel, "debug-level", fmt.Sprintf("possible debug levels: %s",
		strings.Join(value.GetAllLevels(), ", ")))

	return cmd
}
