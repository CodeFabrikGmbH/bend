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
	"net/http"
)

func main() {
	env := createProductionEnvironment()
	defer func() {
		env.Close()
	}()

	http.Handle("/readme/", markdown.FileServer("README.md"))
	http.Handle("/static/", http.FileServer(http.Dir("")))
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/favicon.ico")
	})

	http.Handle("/login", application.LoginPage{Env: env})
	http.Handle("/dashboard/", application.DashboardPage{Env: env})
	http.Handle("/sendRequest/", application.SendRequestService{Env: env})
	http.Handle("/", application.RequestService{Env: env})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func createProductionEnvironment() environment.Environment {
	return environment.Environment{
		LoginPage:         htmlTemplate.LoginPage{},
		DashboardPage:     htmlTemplate.DashBoardPage{},
		RequestRepository: boltDB.CreateRequestRepository(),
		Transport:         _http.Transport{},
		Authentication:    keycloak.New(),
	}
}
