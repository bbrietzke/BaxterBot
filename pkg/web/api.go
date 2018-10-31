package web

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func getStoreValueJSON() (string, http.HandlerFunc) {
	return "/{key:[A-Za-z]+}", func(w http.ResponseWriter, r *http.Request) {
		v := mux.Vars(r)
		logger.Printf("GET: %s", v["key"])

		if value, ok := cache.Get(v["key"]); ok {
			w.Header().Set("Content-Type", "application/json")
			json, _ := json.Marshal(value)
			w.Write([]byte(json))
		} else {
			http.Error(w, "Value not found", http.StatusNotFound)
		}
	}
}

func createStoreValueJSON() (string, http.HandlerFunc) {
	return "/{key:[A-Za-z]+}", func(w http.ResponseWriter, r *http.Request) {
		v := mux.Vars(r)
		logger.Printf("CREATE: %s", v["key"])
		d := json.NewDecoder(r.Body)
		var t interface{}
		err := d.Decode(&t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		cache.Add(v["key"], t)
	}
}
