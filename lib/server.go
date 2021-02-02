package lib

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func IMFetchHandler(db *IM) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(db.Data)
	}
}

func IMUpdateHandler(db *IM) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		fmt.Fprintf(w, "Update id: %v\n", vars["id"])
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(db.Data)
	}
}

func StartServer(db *IM) {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/history", IMFetchHandler(db)).Methods("GET")
	r.HandleFunc("/history/{id}", IMUpdateHandler(db)).Methods("PUT")
	http.Handle("/", r)

	// socketPath := filepath.Join(os.TempDir(), "unixdomain-sample")
	// listener, err := net.Listen("unix", socketPath)

	listener, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		panic(err)
	}

	err = http.Serve(listener, nil)
	if err != nil {
		panic(err)
	}
}
