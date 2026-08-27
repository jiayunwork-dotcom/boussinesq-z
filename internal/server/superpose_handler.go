package server

import (
	"net/http"

	"boussinesq-z/internal/super"
)

type superposeRequest struct {
	Forces []super.Force `json:"forces"`
	Z      float64       `json:"z,omitempty"`
	R      float64       `json:"r,omitempty"`
}

func superposeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, "POST")
		return
	}
	var req superposeRequest
	if !readJSON(w, r, &req) {
		return
	}
	if req.Z == 0 && req.R == 0 {
		result, err := super.Sum(req.Forces)
		if err != nil {
			badRequest(w, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, result)
		return
	}
	result, err := super.CombineAtPoint(req.Forces, req.Z, req.R)
	if err != nil {
		badRequest(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}
