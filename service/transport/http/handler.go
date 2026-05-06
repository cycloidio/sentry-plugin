package http

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/cycloidio/sentry-plugin/service"
)

func Handler(s service.Service) http.Handler {
	r := http.NewServeMux()

	r.Handle("GET /_cy/ping", pingHandler(s))
	r.Handle("POST /_cy/events", eventsHandler(s))
	r.Handle("DELETE /_cy/plugin", deletePluginHandler(s))
	r.Handle("POST /_cy/resync", resyncHandler(s))
	r.Handle("GET /ui/issues", issuesHTMLHandler(s))

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

var issuesTemplate = template.Must(template.New("issues").Parse(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>Sentry Issues</title>
<style>
  body { font-family: sans-serif; margin: 0; padding: 16px; }
  table { border-collapse: collapse; width: 100%; }
  th, td { border: 1px solid #ddd; padding: 8px; text-align: left; font-size: 13px; }
  th { background-color: #f4f4f4; }
  tr:hover { background-color: #f9f9f9; }
  a { color: #0366d6; text-decoration: none; }
  a:hover { text-decoration: underline; }
</style>
</head>
<body>
<h2>Sentry Issues</h2>
<table>
<thead>
<tr>
  <th>Organization</th>
  <th>Project</th>
  <th>Title</th>
  <th>Permalink</th>
  <th>Has Seen</th>
  <th>First Seen</th>
  <th>Last Seen</th>
  <th>User Count</th>
  <th>Level</th>
  <th>Status</th>
  <th>Type</th>
</tr>
</thead>
<tbody>
{{range .}}
<tr>
  <td>{{.OrganizationName}} ({{.OrganizationSlug}})</td>
  <td>{{.ProjectName}} ({{.ProjectSlug}})</td>
  <td>{{.Issue.Title}}</td>
  <td><a href="{{.Issue.Permalink}}" target="_blank">link</a></td>
  <td>{{.Issue.HasSeen}}</td>
  <td>{{.Issue.FirstSeen.Format "2006-01-02 15:04:05"}}</td>
  <td>{{.Issue.LastSeen.Format "2006-01-02 15:04:05"}}</td>
  <td>{{.Issue.UserCount}}</td>
  <td>{{.Issue.Level}}</td>
  <td>{{.Issue.Status}}</td>
  <td>{{.Issue.Type}}</td>
</tr>
{{end}}
</tbody>
</table>
</body>
</html>`))

func issuesHTMLHandler(s service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		issues, err := s.ListAllIssues(r.Context())
		if err != nil {
			http.Error(w, "failed to list issues", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		issuesTemplate.Execute(w, issues)
	}
}
