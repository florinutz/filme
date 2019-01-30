package coll33tx

import (
	"encoding/json"
	"errors"
	"fmt"

	coll33txBusiness "github.com/florinutz/filme/collector/business/coll33tx"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	L33txRootCmd.AddCommand(DetailCmd)
	log.SetFormatter(&log.JSONFormatter{})
}

type dCmdConfigType struct {
	url string
}

var (
	detailsCmdConfig dCmdConfigType

	DetailCmd = &cobra.Command{
		Use:   "detail",
		Short: "Parses the 1337x detail page",

		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("requires the url to be parsed as an argument")
			}
			detailsCmdConfig.url = args[0]
			return nil
		},

		Run: func(cmd *cobra.Command, args []string) {
			log := log.WithField("url", detailsCmdConfig.url)

			detail := coll33txBusiness.NewDetailCollector(log)

			err := detail.Visit(detailsCmdConfig.url)
			if err != nil {
				log.WithError(err).Warn("visit error")
			}

			detail.Wait()

			j, err := json.Marshal(detail.Torrent)
			if err != nil {
				log.Error("could not encode torrent to json")
			}

			fmt.Println(string(j))
		},
	}
)
