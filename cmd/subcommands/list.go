package subcommands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/gocolly/colly"
	"net/url"
	"log"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := colly.NewCollector(
			colly.MaxDepth(1),
		)

		c.OnRequest(func (r *colly.Request) {
			fmt.Println("Visiting url", r.URL)
		})
		c.OnResponse(func(r *colly.Response) {
			fmt.Println("Visited", r.Request.URL)
		})
		c.OnHTML("td.name a:nth-of-type(2)", func(link *colly.HTMLElement) {
			url := fmt.Sprintf("%s://%s%s", rootUrl.Scheme, rootUrl.Host, link.Attr("href"))
			link.Request.Visit(url)
			fmt.Printf("%s (%s)\n", link.Text, url)
		})
		c.OnScraped(func(r *colly.Response) {
			fmt.Println("Finished", r.Request.URL)
		})
		c.OnError(func(r *colly.Response, err error) {
			fmt.Printf("ERROR: %s\n on url %s\nwith esponse: %s\n", err, r.Request.URL, r)
		})

		c.Visit(rootUrl.String())
	},
}

var (
	rootUrl *url.URL
)

func init() {
	rootCmd.AddCommand(listCmd)

	u := listCmd.Flags().StringP("url", "u", "http://1337x.to/popular-movies", "Url to crawl")
	//u := listCmd.Flags().StringP("rootSelector", "u", "http://1337x.to/popular-movies", "Url to crawl")

	x, err := url.Parse(*u)
	if err != nil {
		log.Fatal(err)
	}

	rootUrl = x
}
