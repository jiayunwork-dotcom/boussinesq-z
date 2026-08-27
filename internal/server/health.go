package server

import (
	"net/http"
	"runtime"
	"time"
)

type healthResponse struct {
	Status  string `json:"status"`
	Uptime  int64  `json:"uptime_seconds"`
	GoOS    string `json:"goos"`
	GoArch  string `json:"goarch"`
	Version string `json:"version"`
}

var startedAt = time.Now()

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, "GET")
		return
	}
	writeJSON(w, http.StatusOK, healthResponse{
		Status:  "ok",
		Uptime:  int64(time.Since(startedAt).Seconds()),
		GoOS:    runtime.GOOS,
		GoArch:  runtime.GOARCH,
		Version: "1.0.0",
	})
}
