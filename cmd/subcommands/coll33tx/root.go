package coll33tx

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

var (
	L33txRootCmd = &cobra.Command{
		Use:   "1337x",
		Short: "interactions with 1337x.to",
		Long:  "Handles 1337x.to listings and torrent detail pages.",
	}
)
