package web

import (
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/time/rate"

	"github.com/hashicorp/golang-lru"

	"github.com/gorilla/mux"
)

const (
	requestsPerSeconds int64  = 25
	burstRate          int    = 5
	defaultHTTPPort    string = ":8080"
	defaultCacheSize   int    = 10
)

var (
	logger  *log.Logger
	cache   *lru.Cache
	limiter *rate.Limiter
)

func init() {
	logger = log.New(os.Stdout, "WEB   ", log.LstdFlags|log.Lshortfile)
	cache, _ = lru.New(defaultCacheSize)
	limiter = rate.NewLimiter(rate.Limit(1), burstRate)
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
	Burst
		Maximum number of incoming requests before rate limiting begins.
*/
func Start(options ...Option) error {
	router := mux.NewRouter()
	args := &Options{Port: defaultHTTPPort, RequestPerSecond: requestsPerSeconds, Burst: burstRate, Wait: allowLimitsMW}

	for _, o := range options {
		o(args)
	}

	limiter.SetLimit(rate.Limit(args.RequestPerSecond))

	router.Use(args.Wait)
	router.Use(loggingMW)

	constructAPI(router.PathPrefix("/api").Headers("Content-Type", "application/json").Subrouter())

	srv := &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0" + args.Port,
		WriteTimeout: 20 * time.Second,
		ReadTimeout:  20 * time.Second,
	}

	logger.Printf("starting on %s", args.Port)

	return srv.ListenAndServe()
}
