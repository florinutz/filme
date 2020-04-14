package commands

import (
	"errors"

	"github.com/florinutz/filme/pkg/filme"
	"github.com/spf13/cobra"
)

func Build1337xDetailPageCmd(f *filme.Filme) *cobra.Command {
	var opts struct {
		url                             string
		justMagnet                      bool
		json                            bool
		delay, randomDelay, parallelism int
		userAgent                       string
	}

	cmd := &cobra.Command{
		Use:   "1337x_detail <url>",
		Short: "Parses the 1337x detail page. Set OMDB_API_KEY if you want imdb info",

		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("requires the url to be parsed as an argument")
			}
			opts.url = args[0]

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			return f.Visit1337xDetailPage(opts.url, opts.justMagnet, opts.json,
				opts.delay, opts.randomDelay, opts.parallelism, opts.userAgent)
		},
	}

	cmd.Flags().BoolVarP(&opts.justMagnet, "magnet", "m", false, "only show magnet")
	cmd.Flags().BoolVarP(&opts.json, "json", "j", false, "encode torrent output to json")
	cmd.Flags().IntVar(&opts.delay, "reqs-delay", 3, "requests delay")
	cmd.Flags().IntVar(&opts.randomDelay, "reqs-random-delay", 3, "requests random delay")
	cmd.Flags().IntVar(&opts.parallelism, "reqs-parallelism", 4, "requests parallelism")
	cmd.Flags().StringVar(&opts.userAgent, "user-agent", "", "requests user agent. Leave this empty for random browser UAs")

	return cmd
}
