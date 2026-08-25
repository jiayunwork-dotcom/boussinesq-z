package server

import (
	"encoding/json"
	"net/http"

	"boussinesq-z/internal/circular"
)

func Routes(store *circular.Store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/api/stress", stressHandler)
	mux.HandleFunc("/api/superpose", superposeHandler)
	mux.HandleFunc("/api/circular", circularHandler(store))
	mux.HandleFunc("/api/table-snapshot", tableSnapshotHandler(store))
	mux.HandleFunc("/api/version", versionHandler)
	return mux
}

func requireMethod(next http.HandlerFunc, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			writeError(w, http.StatusMethodNotAllowed, "method must be "+method)
			return
		}
		next(w, r)
	}
}

func readJSON(w http.ResponseWriter, r *http.Request, target interface{}) bool {
	if r.Body == nil {
		writeError(w, http.StatusBadRequest, "missing request body")
		return false
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return false
	}
	return true
}
