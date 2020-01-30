package filme

import (
	"fmt"
	"net/http"
	"os"
)

func (f *Filme) HttpHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(f.Out, "received a request.\n\t UA: %s\n", req.UserAgent())
		target := os.Getenv("TARGET")
		if target == "" {
			target = "World"
		}
		fmt.Fprintf(w, "Hello %s!\n", target)
	}
}
