package main

import (
	"code-fabrik.com/bend/application"
	"code-fabrik.com/bend/domain/config"
	"code-fabrik.com/bend/domain/request"
	"code-fabrik.com/bend/infrastructure/boltDB"
	httptransport "code-fabrik.com/bend/infrastructure/http"
	"code-fabrik.com/bend/infrastructure/httpHandler"
	"code-fabrik.com/bend/infrastructure/jwt/keycloak"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/google/uuid"
	"net/http"
)

func main() {
	requestRepository, configRepository, transport, db := createProductionEnvironment()
	defer func() {
		_ = db.Close()
	}()

	migrate(configRepository)

	keycloakService := keycloak.New()
	configService := application.ConfigService{
		ConfigRepository: configRepository,
	}

	requestService := application.RequestService{
		RequestRepository: requestRepository,
		ConfigRepository:  configRepository,
		Transport:         transport,
	}

	dashboardService := application.DashboardService{
		RequestRepository: requestRepository,
	}

	http.Handle("/static/", http.FileServer(http.Dir("")))
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/favicon.ico")
	})

	http.Handle("/readme/", httpHandler.ReadMePage{MarkdownFile: "README.md"})
	http.Handle("/login", httpHandler.LoginPage{KeyCloakService: keycloakService})
	http.Handle("/dashboard/", httpHandler.DashboardPage{KeyCloakService: keycloakService, DashboardService: dashboardService})
	http.Handle("/configs/", httpHandler.ConfigPage{KeyCloakService: keycloakService, ConfigService: configService})

	http.Handle("/api/configs/", httpHandler.ConfigAPI{KeyCloakService: keycloakService, ConfigService: configService})
	http.Handle("/api/requests/", httpHandler.RequestAPI{RequestService: requestService})

	http.Handle("/", httpHandler.TrackRequest{RequestService: requestService})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
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
		fmt.Println("adding identifier and deleting all entries with URL keys")
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
