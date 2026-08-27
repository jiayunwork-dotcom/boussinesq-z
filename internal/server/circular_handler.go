package server

import (
	"net/http"

	"boussinesq-z/internal/circular"
)

type circularRequest struct {
	Q float64 `json:"q"`
	A float64 `json:"a"`
	Z float64 `json:"z"`
	R float64 `json:"r"`
}

func circularHandler(store *circular.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			methodNotAllowed(w, "POST")
			return
		}
		var req circularRequest
		if !readJSON(w, r, &req) {
			return
		}
		load := circular.Load{Q: req.Q, A: req.A, Z: req.Z, R: req.R}
		result, err := store.Evaluate(load)
		if err != nil {
			badRequest(w, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, result)
	}
}

func tableSnapshotHandler(store *circular.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			methodNotAllowed(w, "GET")
			return
		}
		snapshot, err := store.Snapshot()
		if err != nil {
			badRequest(w, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, snapshot)
	}
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, "GET")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"name":    "boussinesq-z",
		"version": "1.0.0",
	})
}
