package coll33tx

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	coll33txBusiness "gitlab.com/phlo/filme/collector/business/coll33tx"
)

func init() {
	Coll33txCmd.AddCommand(DetailCmd)
	log.SetFormatter(&log.JSONFormatter{})
}

var (
	DetailCmd = &cobra.Command{
		Use:   "detail",
		Short: "Parses the 1337x detail page",
		Long:  "Handles 1337x.to detail pages",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("requires the url to be parsed as an argument")
			}
			Url = args[0]
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			log := log.WithField("url", Url)

			detail := coll33txBusiness.NewDetailCollector(log)

			err := detail.Visit(Url)
			if err != nil {
				log.WithError(err).Warn("visit error")
			}

			detail.Wait()
		},
	}
)
