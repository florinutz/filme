package coll33tx

import (
	coll33txBusiness "github.com/florinutz/filme/collector/business/coll33tx"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
