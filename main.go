package main

import (
	"code-fabrik.com/bend/application"
	"code-fabrik.com/bend/domain/config"
	"code-fabrik.com/bend/domain/request"
	"code-fabrik.com/bend/infrastructure/boltDB"
	"code-fabrik.com/bend/infrastructure/env"
	"code-fabrik.com/bend/infrastructure/htmlTemplate"
	httptransport "code-fabrik.com/bend/infrastructure/http"
	"code-fabrik.com/bend/infrastructure/httpHandler"
	"code-fabrik.com/bend/infrastructure/jwt/keycloak"
	"context"
	"errors"
	"github.com/google/uuid"
	bolt "go.etcd.io/bbolt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := htmlTemplate.Load(resourcesFS, "resources/*.html"); err != nil {
		slog.Error("failed to load templates", "err", err)
		os.Exit(1)
	}

	requestRepository, configRepository, transport, db := createProductionEnvironment()
	defer func() {
		_ = db.Close()
	}()

	migrate(configRepository)

	// One-off backfill of the per-path request counts used by the dashboard.
	// Runs in the background so it never blocks start-up or request tracking; it
	// is a no-op after the first successful run.
	if backfiller, ok := requestRepository.(interface{ BackfillPathCounts() error }); ok {
		go func() {
			if err := backfiller.BackfillPathCounts(); err != nil {
				slog.Error("path count backfill failed", "err", err)
			} else {
				slog.Info("path count backfill complete")
			}
		}()
	}

	keycloakService := keycloak.New()
	eventHub := application.NewEventHub()

	configService := application.ConfigService{
		ConfigRepository: configRepository,
	}

	requestService := application.RequestService{
		RequestRepository: requestRepository,
		ConfigRepository:  configRepository,
		Transport:         transport,
		Hub:               eventHub,
	}

	dashboardService := application.DashboardService{
		RequestRepository: requestRepository,
	}

	mux := http.NewServeMux()

	mux.Handle("/static/", http.FileServer(http.FS(staticFS)))
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, staticFS, "static/favicon.ico")
	})

	mux.Handle("/readme/", httpHandler.ReadMePage{Markdown: readmeMarkdown})
	mux.Handle("/login", httpHandler.LoginPage{KeyCloakService: keycloakService})
	mux.Handle("/dashboard/", httpHandler.DashboardPage{KeyCloakService: keycloakService, DashboardService: dashboardService})
	mux.Handle("/configs/", httpHandler.ConfigPage{KeyCloakService: keycloakService, ConfigService: configService})

	mux.Handle("/api/configs/", httpHandler.ConfigAPI{KeyCloakService: keycloakService, ConfigService: configService})
	mux.Handle("/api/requests/", httpHandler.RequestAPI{KeyCloakService: keycloakService, RequestService: requestService})
	mux.Handle("/api/request-list", httpHandler.RequestListAPI{KeyCloakService: keycloakService, RequestService: requestService})
	mux.Handle("/api/events", httpHandler.EventsAPI{KeyCloakService: keycloakService, Hub: eventHub})

	// The bare root serves the admin UI (config page); every other path is the
	// request catcher. "/{$}" matches only "/", so it takes precedence over the
	// catch-all "/" below without shadowing tracked paths.
	mux.Handle("/{$}", http.RedirectHandler("/configs/", http.StatusFound))
	mux.Handle("/", httpHandler.TrackRequest{RequestService: requestService})

	server := &http.Server{
		Addr:              env.LISTEN_ADDR,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	go func() {
		slog.Info("bend listening", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	slog.Info("shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("graceful shutdown failed", "err", err)
	}
}

func migrate(configRepository config.Repository) {
	configs := configRepository.FindAll()
	addIdentifier := false
	for _, configItem := range configs {
		if configItem.Id == uuid.Nil {
			addIdentifier = true
		}
	}
	if addIdentifier {
		slog.Info("adding identifier and deleting all entries with URL keys")
		_ = configRepository.DeleteAll()
		for _, configItem := range configs {
			configItem.Id = uuid.New()
			_ = configRepository.Save(configItem)
		}
	}
}

func createProductionEnvironment() (request.Repository, config.Repository, httptransport.Transport, *bolt.DB) {
	db, err := bolt.Open("db/my.db", 0600, nil)
	if err != nil {
		panic(err)
	}

	return boltDB.RequestRepository{DB: db}, boltDB.ConfigRepository{DB: db}, httptransport.Transport{}, db
}
