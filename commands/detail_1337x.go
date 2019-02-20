package commands

import (
	"errors"

	"github.com/florinutz/filme/pkg/filme"
	"github.com/spf13/cobra"
)

func Build1337xParseDetailPageCmd(f *filme.Filme) *cobra.Command {
	var opts struct {
		url        string
		justMagnet bool
		json       bool
	}

	cmd := &cobra.Command{
		Use:   "detail",
		Short: "Parses the 1337x detail page",

		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("requires the url to be parsed as an argument")
			}
			opts.url = args[0]

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			return f.VisitDetailPage(opts.url, opts.justMagnet, opts.json)
		},
	}

	cmd.Flags().BoolVarP(&opts.justMagnet, "magnet", "m", false,
		"only show magnet")
	cmd.Flags().BoolVarP(&opts.json, "json", "j", false,
		"encode torrent output to json")

	return cmd
}
