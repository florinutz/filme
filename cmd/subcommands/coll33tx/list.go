package coll33tx

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	coll33txBusiness "gitlab.com/phlo/filme/collector/business/coll33tx"
)

func init() {
	ListCmd.Flags().BoolVarP(&withDetails, "crawl-details", "d", false, "follows every link in the list and fetches more data")

	Coll33txCmd.AddCommand(ListCmd)

	log.SetFormatter(&log.JSONFormatter{})
}

var (
	withDetails bool

	ListCmd = &cobra.Command{
		Use:   "list",
		Short: "Parses 1337x listings",
		Long:  "Handles 1337x.to listings.",

		Run: func(cmd *cobra.Command, args []string) {
			list := coll33txBusiness.NewListCollector(withDetails, log.WithField("url", Url))

			err := list.Visit(Url)
			if err != nil {
				log.WithError(err).WithField("url", Url).Warn("visit error")
			}

			list.Wait()
		},
	}
)
