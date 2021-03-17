package main

import (
	"code-fabrik.com/bend/application"
	"code-fabrik.com/bend/domain/environment"
	"code-fabrik.com/bend/infrastructure/boltDB"
	httptransport "code-fabrik.com/bend/infrastructure/http"
	"code-fabrik.com/bend/infrastructure/httpHandler"
	"code-fabrik.com/bend/infrastructure/jwt/keycloak"
	"fmt"
	"github.com/boltdb/bolt"
	"net/http"
)

func main() {
	env, db := createProductionEnvironment()
	defer func() {
		_ = db.Close()
	}()

	keycloakService := keycloak.New()
	configService := application.ConfigService{Env: env}
	requestService := application.RequestService{Env: env}
	dashboardService := application.DashboardService{Env: env}

	http.Handle("/static/", http.FileServer(http.Dir("")))
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/favicon.ico")
	})

	http.Handle("/readme/", httpHandler.ReadMePage{MarkdownFile: "README.md"})
	http.Handle("/login", httpHandler.LoginPage{KeyCloakService: keycloakService})
	http.Handle("/dashboard/", httpHandler.DashboardPage{KeyCloakService: keycloakService, DashboardService: dashboardService})
	http.Handle("/config/", httpHandler.ConfigPage{KeyCloakService: keycloakService, ConfigService: configService})

	http.Handle("/delete/", httpHandler.Deletion{RequestService: requestService})
	http.Handle("/sendRequest/", httpHandler.SendRequest{SendRequestService: requestService})

	http.Handle("/", httpHandler.TrackableRequest{RequestService: requestService})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func createProductionEnvironment() (environment.Environment, *bolt.DB) {
	db, err := bolt.Open("db/my.db", 0600, nil)
	if err != nil {
		panic(err)
	}

	return environment.Environment{
		RequestRepository: boltDB.RequestRepository{DB: db},
		ConfigRepository:  boltDB.ConfigRepository{DB: db},
		Transport:         httptransport.Transport{},
	}, db
}
