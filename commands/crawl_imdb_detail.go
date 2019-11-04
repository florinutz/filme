package commands

import (
	"errors"

	"github.com/florinutz/filme/pkg/filme"
	"github.com/spf13/cobra"
)

func BuildImdbDetailPageCmd(f *filme.Filme) *cobra.Command {
	var opts struct {
		url  string
		json bool
	}

	cmd := &cobra.Command{
		Use:   "imdb",
		Short: "Parses the imdb detail page",

		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("requires the url to be parsed as an argument")
			}
			opts.url = args[0]

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			return f.VisitImdbDetailPage(opts.url, opts.json)
		},
	}

	cmd.Flags().BoolVarP(&opts.json, "json", "j", false,
		"encode output to json")

	return cmd
}
