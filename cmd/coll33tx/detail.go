package coll33tx

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/eefret/gomdb"

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
			logger := &log.Logger{
				Out:          os.Stderr,
				Formatter:    &log.TextFormatter{},
				Level:        log.DebugLevel,
				ReportCaller: true,
			}

			log := logger.WithFields(log.Fields{
				"url":  detailsCmdConfig.url,
				"type": "detail_page",
			})

			detail := coll33tx.NewDetailsPageCollector(torrentFound, log)

			if err := detail.Visit(detailsCmdConfig.url); err != nil {
				log.WithError(err).Fatal("visit error")
			}

			detail.Wait()
		},
	}
)

func torrentFound(torrent coll33tx.Torrent) {
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
		torrent.Seeders,
		torrent.Leechers,
	)

	if omdbApiKey, ok := os.LookupEnv("OMDB_API_KEY"); ok {
		gomdbApi := gomdb.Init(omdbApiKey)
		query := &gomdb.QueryData{Title: torrent.FilmCleanTitle, SearchType: gomdb.MovieSearch}
		res, err := gomdbApi.Search(query)
		if err != nil {
			log.WithError(err).WithField("title", torrent.Title).Fatal("omdb lookup failed")
		}
		fmt.Println(res.Search)
	} else {
		log.Warn("no omdb api key")
	}

}
