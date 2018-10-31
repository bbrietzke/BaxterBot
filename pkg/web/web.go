package web

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/golang-lru"

	"github.com/gorilla/mux"
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

// Start creates and starts an HTTP server on port 8080, or a different port
// if one is provided as an option.
//
// Options that can be included are:
// * Port
//
func Start(options ...Option) error {
	args := &Options{Port: ":8080"}

	router.Use(loggingMW)
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

func loggingMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

/*
curl -X GET -H "Content-type: application/json" -H "Accept: application/json"  "http://localhost:8080/api/fred"

curl --header "Content-Type: application/json" --request POST --data '{"username":"xyz","password":"xyz"}' http://localhost:8080/api/fred
*/
