package commands

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"os"
	"strconv"

	"github.com/florinutz/filme/infra/proto"

	"github.com/asim/go-micro/v3"

	"github.com/florinutz/filme/pkg/filme"
	"github.com/spf13/cobra"
)

type Crawler struct{}

func (c Crawler) Search(ctx context.Context, request *proto.SearchRequest, response *proto.SearchResponse) error {
	log.Infof("Received a Search request for term \"%s\"", request.Term)
	response.Term = fmt.Sprintf("ecce %s", request.Term)
	return nil
}

// BuildServeCmd mirrors the 1337x search cmd
func BuildServeCmd(f *filme.Filme) (cmd *cobra.Command) {
	var opts struct {
		port int
		nats string
	}

	cmd = &cobra.Command{
		Use:   "serve",
		Short: "serve search over http endpoints",

		RunE: func(cmd *cobra.Command, args []string) error {
			go func() {
				for {
					grpc.DialContext(context.TODO(), "127.0.0.1:9091")
					time.Sleep(time.Second)
				}
			}()

			service := micro.NewService(
				micro.Name("crawl"),
				micro.Version("latest"),
			)
			service.Init(micro.AfterStart(func() error {
				log.Infof("service listening on %s!",
					service.Options().Server.Options().Address,
				)
				return nil
			}))

			if err := proto.RegisterTorrentsHandler(service.Server(), new(Crawler)); err != nil {
				return err
			}

			if err := service.Run(); err != nil {
				log.Fatal(err)
			}

			return nil
		},
	}

	var (
		portStr string
		err     error
	)

	if portStr = os.Getenv("PORT"); portStr == "" {
		portStr = "8080"
	}
	if opts.port, err = strconv.Atoi(portStr); err != nil {
		f.Log.Fatalf("port not int")
	}

	opts.nats = os.Getenv("NATS_DSN")

	cmd.Flags().IntVarP(&opts.port, "port", "p", opts.port, "listen port")
	cmd.Flags().StringVarP(&opts.nats, "nats-host", "n", opts.nats, "listen port")

	return
}
