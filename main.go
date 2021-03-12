package main

import (
	"code-fabrik.com/bend/application"
	"code-fabrik.com/bend/domain/environment"
	"code-fabrik.com/bend/infrastructure/_http"
	"code-fabrik.com/bend/infrastructure/boltDB"
	"code-fabrik.com/bend/infrastructure/htmlTemplate"
	"code-fabrik.com/bend/infrastructure/jwt/keycloak"
	"code-fabrik.com/bend/infrastructure/markdown"
	"fmt"
	"github.com/boltdb/bolt"
	"net/http"
)

func main() {
	env, db := createProductionEnvironment()
	defer func() {
		_ = db.Close()
	}()

	http.Handle("/readme/", markdown.FileServer("README.md"))
	http.Handle("/static/", http.FileServer(http.Dir("")))
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/favicon.ico")
	})

	http.Handle("/login", application.LoginPage{Env: env})
	http.Handle("/dashboard/", application.DashboardPage{Env: env})
	http.Handle("/delete/", application.DeletionService{Env: env})
	http.Handle("/sendRequest/", application.SendRequestService{Env: env})

	http.Handle("/", application.RequestService{Env: env})

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
		LoginPage:         htmlTemplate.LoginPagePresenter{},
		DashboardPage:     htmlTemplate.DashBoardPage{},
		RequestRepository: boltDB.RequestRepository{DB: db},
		ConfigRepository:  boltDB.ConfigRepository{DB: db},
		Transport:         _http.Transport{},
		Authentication:    keycloak.New(),
	}, db
}
