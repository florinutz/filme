package coll33tx

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	Coll33txCmd.PersistentFlags().StringVarP(
		&Url,
		"url",
		"u",
		"",
		"the url to crawl. which can either be a listing page or a details page",
	)
	Coll33txCmd.MarkFlagRequired("url")
	log.SetFormatter(&log.JSONFormatter{})
}

var (
	Url string

	Coll33txCmd = &cobra.Command{
		Use:   "1337x",
		Short: "interactions with 1337x.to",
		Long:  "Handles 1337x.to listings and torrent detail pages.",
	}
)
