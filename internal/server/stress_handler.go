package server

import (
	"net/http"

	"boussinesq-z/internal/point"
)

type stressRequest struct {
	P       float64 `json:"P"`
	Z       float64 `json:"z"`
	R       float64 `json:"r"`
	Poisson float64 `json:"poisson,omitempty"`
}

func stressHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, "POST")
		return
	}
	var req stressRequest
	if !readJSON(w, r, &req) {
		return
	}
	poisson := req.Poisson
	if poisson == 0 {
		poisson = 0.3
	}
	if err := point.ValidatePoisson(poisson); err != nil {
		badRequest(w, err.Error())
		return
	}
	load := point.Load{P: req.P, Z: req.Z, R: req.R}
	result, err := point.Evaluate(load, poisson)
	if err != nil {
		badRequest(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}
