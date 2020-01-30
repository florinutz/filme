package commands

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/florinutz/filme/pkg/filme"
	"github.com/spf13/cobra"
)

// BuildSearchCmd mirrors the 1337x search cmd
func BuildServeCmd(f *filme.Filme) (cmd *cobra.Command) {
	var opts struct {
		port int
	}

	cmd = &cobra.Command{
		Use:   "serve",
		Short: "serve search over http endpoints",

		RunE: func(cmd *cobra.Command, args []string) error {
			f.Log.Println("Hello filme sample started.")

			http.HandleFunc("/", f.HttpHandler())

			port := os.Getenv("PORT")
			if port == "" {
				port = "8080"
			}

			return http.ListenAndServe(fmt.Sprintf(":%d", opts.port), nil)
		},
	}

	var (
		portStr string
		port    int
		err     error
	)

	if portStr = os.Getenv("PORT"); portStr == "" {
		portStr = "8080"
	}

	if port, err = strconv.Atoi(portStr); err != nil {
		f.Log.Fatalf("port not int")
	}

	cmd.Flags().IntVarP(&opts.port, "port", "p", port, "listen port")

	return
}
