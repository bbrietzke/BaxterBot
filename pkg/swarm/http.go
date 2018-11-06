package swarm

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/hashicorp/raft"
)

type registration struct {
	Address string
	Name    string
}

// SetupHTTP takes an http.ServeMux and adds a few paths to it.
func SetupHTTP(mux *mux.Router) {
	mux.Handle(handleJoin()).Methods("POST")
	mux.Handle(handleLeave()).Methods("POST")
}

func handleLeave() (string, http.HandlerFunc) {
	return "/leave/{node}", func(w http.ResponseWriter, r *http.Request) {
		if !isLeader() {
			logger.Println("Redirecting to leader:", leaderRedirect)
			http.Redirect(w, r, leaderRedirect, http.StatusPermanentRedirect)
			return
		}
	}
}

func handleJoin() (string, http.HandlerFunc) {
	return "/join", func(w http.ResponseWriter, r *http.Request) {
		if !isLeader() {
			logger.Println("Redirecting to leader:", leaderRedirect)
			http.Redirect(w, r, leaderRedirect, http.StatusPermanentRedirect)
			return
		}

		reg := registration{}
		if err := json.NewDecoder(r.Body).Decode(&reg); err != nil {
			logger.Println("bad registration request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		cfg := swarmer.GetConfiguration()
		if err := cfg.Error(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logger.Println(err)
			return
		}

		for _, srv := range cfg.Configuration().Servers {
			if srv.ID == raft.ServerID(reg.Name) || srv.Address == raft.ServerAddress(reg.Address) {
				if srv.ID == raft.ServerID(reg.Name) && srv.Address == raft.ServerAddress(reg.Address) {
					logger.Println(srv.ID, " is already a member")
					return
				}

				if err := swarmer.RemoveServer(srv.ID, 0, 0).Error(); err != nil {
					logger.Println(srv.ID, " is already a member and could not be removed")
					logger.Println(err)
					return
				}

			}
		}

		if err := swarmer.AddVoter(raft.ServerID(reg.Name), raft.ServerAddress(reg.Address), 0, time.Second*5); err != nil {
			logger.Println("node", reg.Name, "has joined.")
		}
	}
}
