package web

import (
	"encoding/json"
	"net/http"

	"github.com/bbrietzke/BaxterBot/pkg/swarm"

	"github.com/gorilla/mux"
)

func constructAPI(r *mux.Router) {
	r.Handle(createStoreValueJSON()).Methods("POST")
	r.Handle(getStoreValueJSON()).Methods("GET")
	r.Handle(deleteStoreValueJSON()).Methods("DELETE")
}

func getStoreValueJSON() (string, http.HandlerFunc) {
	return "/{key:[A-Za-z]+}", func(w http.ResponseWriter, r *http.Request) {
		v := mux.Vars(r)

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

		d := json.NewDecoder(r.Body)
		var t interface{}
		err := d.Decode(&t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		cache.Add(v["key"], t)
		swarm.CreateKeyValueEntry(v["key"], t)
		w.WriteHeader(http.StatusOK)
	}
}

func deleteStoreValueJSON() (string, http.HandlerFunc) {
	return "/{key:[A-Za-z]+}", func(w http.ResponseWriter, r *http.Request) {
		v := mux.Vars(r)
		cache.Remove(v["key"])
		w.WriteHeader(http.StatusOK)
	}
}
