package web

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/golang-lru"

	"github.com/gorilla/mux"
)

const (
	requestsPerSeconds int64  = 10
	burstRate          int    = 2
	defaultHTTPPort    string = ":8080"
)

var (
	router *mux.Router
	logger *log.Logger
	cache  *lru.Cache
)

func init() {
	router = mux.NewRouter()
	logger = log.New(os.Stdout, "WEB  ", log.LstdFlags|log.Lshortfile)
	cache, _ = lru.New(8)
}

/*
Start creates and starts an HTTP server on port 8080, or a different port
if one is provided as an option.

Options that can be included are:
	Port
		Sets the port that the http.Server will run on.
	RequestsPerSecond
		Sets the rate limiter to only allow a set amount of HTTP requests per second.
		Automatically sets the Burst to be 10% of the RequestsPerSecond value.
*/
func Start(options ...Option) error {
	args := &Options{Port: defaultHTTPPort, RequestPerSecond: requestsPerSeconds, Burst: burstRate}

	for _, o := range options {
		o(args)
	}

	router.Use(loggingMW)
	router.Use(rateLimitingMW(args.RequestPerSecond, args.Burst))
	api := router.PathPrefix("/api").Headers("Content-Type", "application/json").Subrouter()

	api.Handle(createStoreValueJSON()).Methods("POST")
	api.Handle(getStoreValueJSON()).Methods("GET")

	srv := &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0" + args.Port,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	logger.Printf("starting on %s", args.Port)

	return srv.ListenAndServe()
}

// cat body.txt | vegeta attack -duration 10s  | tee /tmp/report.bin | vegeta report -type=text && cat /tmp/report.bin | vegeta.plot > /tmp/open.html && open /tmp/page.html
