package coll33tx

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/florinutz/filme/collector/coll33tx"

	coll33txBusiness "github.com/florinutz/filme/collector/business/coll33tx"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	L33txRootCmd.AddCommand(DetailCmd)
	log.SetFormatter(&log.JSONFormatter{})
}

type dCmdConfigType struct {
	url        string
	justMagnet bool
}

func init() {
	DetailCmd.Flags().BoolVarP(&detailsCmdConfig.justMagnet, "magnet", "m", false,
		"only show magnet")
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
			log.SetReportCaller(true)
			log := log.WithField("url", detailsCmdConfig.url)

			detail := coll33txBusiness.NewDetailCollector(log)

			err := detail.Visit(detailsCmdConfig.url)
			if err != nil {
				log.WithError(err).Warn("visit error")
			}

			detail.Wait()

			if err := displayTorrent(detail.Torrent); err != nil {
				log.WithError(err).Fatal("could not display torrent")
			}
		},
	}
)

func displayTorrent(torrent coll33tx.L33tTorrent) error {
	if detailsCmdConfig.justMagnet {
		fmt.Println(torrent.Magnet)
		return nil
	}

	j, err := json.Marshal(torrent)
	if err != nil {
		return err
	}
	fmt.Println(string(j))

	return nil
}
