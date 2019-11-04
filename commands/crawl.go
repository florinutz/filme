package commands

import (
	"github.com/florinutz/filme/pkg/filme"
	"github.com/spf13/cobra"
)

func BuildCrawlCmd(f *filme.Filme) *cobra.Command {
	return &cobra.Command{
		Use:   "crawl",
		Short: "Groups crawl commands",
	}
}
