package filme

import (
	"fmt"
	"net/http"
	"os"
)

func (f *Filme) HttpHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		f.Log.Println("Hello world received a request.")
		target := os.Getenv("TARGET")
		if target == "" {
			target = "World"
		}
		fmt.Fprintf(w, "Hello %s!\n", target)

	}
}
