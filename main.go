package main

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"io"
	"os"

	"github.com/fastly/compute-sdk-go/fsthttp"
	"github.com/mealies/wasmScrabbleScorer/internal/handlers"
)

//go:embed static/index.html
var indexPage []byte

func main() {
	// Log service version
	fmt.Println("FASTLY_SERVICE_VERSION:", os.Getenv("FASTLY_SERVICE_VERSION"))

	fsthttp.ServeFunc(func(ctx context.Context, w fsthttp.ResponseWriter, r *fsthttp.Request) {
		switch r.URL.Path {
		case "/api/start":
			handlers.HandleStart(w, r)
		case "/api/score":
			handlers.HandleScore(w, r)
		case "/api/end-round":
			handlers.HandleEndRound(w, r)
		case "/api/finish":
			handlers.HandleFinish(w, r)
		case "/api/status":
			handlers.HandleStatus(w, r)
		case "/":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			io.Copy(w, bytes.NewReader(indexPage))
		default:
			// Catch all other requests and return a 404.
			w.WriteHeader(fsthttp.StatusNotFound)
			fmt.Fprintf(w, "The page you requested could not be found\n")
		}
	})
}
