package main

import (
	"code-fabrik.com/bend/application"
	"code-fabrik.com/bend/domain/dashboard"
	"code-fabrik.com/bend/domain/request"
	"code-fabrik.com/bend/infrastructure/_http"
	"code-fabrik.com/bend/infrastructure/boltDB"
	"code-fabrik.com/bend/infrastructure/htmlTemplate"
	"code-fabrik.com/bend/infrastructure/markdown"
	"fmt"
	"net/http"
)

func main() {
	requestRepository, dashboardPresenter, transport := createProductionEnvironment()
	defer func() {
		requestRepository.Close()
	}()

	http.Handle("/readme/", markdown.FileServer("README.md"))
	http.Handle("/static/", http.FileServer(http.Dir("")))
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/favicon.ico")
	})

	http.Handle("/dashboard/", application.DashboardService{
		RequestRepository:  requestRepository,
		DashboardPresenter: dashboardPresenter,
	})

	http.Handle("/sendRequest/", application.SendRequestService{
		RequestRepository: requestRepository,
		TransportService:  transport,
	})

	http.Handle("/", application.RequestService{
		RequestRepository: requestRepository,
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func createProductionEnvironment() (request.Repository, dashboard.Presenter, request.Transport) {
	return boltDB.Create(), htmlTemplate.DashBoardPresenter{}, _http.Transport{}
}
