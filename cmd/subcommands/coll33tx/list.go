package coll33tx

import (
	"fmt"
	"net/url"
	"strings"

	coll33txBusiness "github.com/florinutz/filme/collector/business/coll33tx"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	ListCmd.Flags().BoolVarP(&listCmdConfig.withDetails, "crawl-details", "d", false, "follows every link in the list and fetches more data")
	L33txRootCmd.AddCommand(ListCmd)
	log.SetFormatter(&log.JSONFormatter{})
}

type lCmdConfigType struct {
	withDetails bool
	url         string
}

var (
	listCmdConfig lCmdConfigType

	ListCmd = &cobra.Command{
		Use:   "list",
		Short: "Parses 1337x listings",

		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				if strings.Contains(args[0], "//1337x.to") {
					listCmdConfig.url = args[0]
				} else { // search
					params := url.Values{"search": {strings.Join(args, " ")}}
					listCmdConfig.url = fmt.Sprintf("https://1337x.to/srch?%s", params.Encode())
				}
			} else {
				listCmdConfig.url = "https://1337x.to/popular-movies"
			}
			return nil
		},

		Run: func(cmd *cobra.Command, args []string) {
			log := log.WithField("url", listCmdConfig.url)
			list := coll33txBusiness.NewListCollector(listCmdConfig.withDetails, log)
			err := list.Visit(listCmdConfig.url)
			if err != nil {
				log.WithError(err).Warn("visit error")
			}
			list.Wait()
		},
	}
)
