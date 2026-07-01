package httpHandler

import (
	"code-fabrik.com/bend/application"
	"code-fabrik.com/bend/infrastructure/jwt/keycloak"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// EventsAPI streams newly tracked requests to the browser via Server-Sent
// Events, driving the dashboard's live tail.
type EventsAPI struct {
	KeyCloakService *keycloak.Service
	Hub             *application.EventHub
}

func (e EventsAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			slog.Error("panic in EventsAPI", "recover", rec)
		}
	}()

	if _, err := e.KeyCloakService.Authenticate(w, r); err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// The server sets a global WriteTimeout; SSE connections are long-lived, so
	// clear the write deadline for this connection only.
	rc := http.NewResponseController(w)
	_ = rc.SetWriteDeadline(time.Time{})

	events := e.Hub.Subscribe()
	defer e.Hub.Unsubscribe(events)

	// open the stream and prompt proxies to start forwarding
	_, _ = fmt.Fprint(w, ": connected\n\n")
	_ = rc.Flush()

	// periodic heartbeat keeps intermediaries from closing an idle connection
	ticker := time.NewTicker(25 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			if _, err := fmt.Fprint(w, ": ping\n\n"); err != nil {
				return
			}
			_ = rc.Flush()
		case event := <-events:
			data, err := json.Marshal(event)
			if err != nil {
				continue
			}
			if _, err := fmt.Fprintf(w, "data: %s\n\n", data); err != nil {
				return
			}
			_ = rc.Flush()
		}
	}
}
