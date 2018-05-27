package subcommands

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.com/phlo/filme/collector/l33tx"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		l33tx.ListCollector.Visit("https://1337x.to/popular-movies")
		l33tx.ListCollector.Wait()

		torrents := l33tx.Torrents
		fmt.Printf("Found %d torrents:\n\n", len(torrents))
		for _, torrent := range torrents {
			fmt.Printf("%s\n\t%s\n\t%s\n\n", torrent.Title, torrent.FoundOn, torrent.Magnet)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
