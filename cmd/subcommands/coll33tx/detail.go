package coll33tx

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/florinutz/filme/collector/coll33tx"

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
	json       bool
}

func init() {
	DetailCmd.Flags().BoolVarP(&detailsCmdConfig.justMagnet, "magnet", "m", false,
		"only show magnet")
	DetailCmd.Flags().BoolVarP(&detailsCmdConfig.json, "json", "j", false,
		"encode torrent output to json")
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

			detail := coll33tx.NewDetailsPageCollector(torrentFound)

			err := detail.Visit(detailsCmdConfig.url)
			if err != nil {
				log.WithError(err).Warn("visit error")
			}

			detail.Wait()
		},
	}
)

func torrentFound(torrent coll33tx.L33tTorrent) {
	log.WithField("torrent", torrent).Debug("torrent found on detail page")

	if detailsCmdConfig.justMagnet {
		fmt.Println(torrent.Magnet)
		return
	}

	if detailsCmdConfig.json {
		j, err := json.Marshal(torrent)
		if err != nil {
			log.WithError(err).Fatal("error encoding to json")
		}
		fmt.Println(string(j))
		return
	}

	fmt.Printf(`%s

magnet: %s

seeders: %d
leechers: %d`,
		strings.Trim(torrent.Title, " "),
		torrent.Magnet,
		torrent.Seeds,
		torrent.Leeches,
	)
}
