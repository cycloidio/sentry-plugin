package http

import (
	"encoding/json"
	"net/http"

	"github.com/cycloidio/sentry-plugin/service"
)

func Handler(s service.Service) http.Handler {
	r := http.NewServeMux()

	r.Handle("GET /_cy/ping", pingHandler(s))
	r.Handle("POST /_cy/events", eventsHandler(s))
	r.Handle("DELETE /_cy/plugin", deletePluginHandler(s))
	r.Handle("POST /_cy/resync", resyncHandler(s))

	return r
}

func pingHandler(s service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"request": "ping"})
	}
}

func eventsHandler(s service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"request": "events"})
	}
}

func deletePluginHandler(s service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"request": "plugin"})
	}
}

func resyncHandler(s service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"request": "resync"})
	}
}
