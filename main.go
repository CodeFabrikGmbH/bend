package main

import (
	"code-fabrik.com/bend/application"
	"code-fabrik.com/bend/infrastructure/_http"
	"code-fabrik.com/bend/infrastructure/boltDB"
	"code-fabrik.com/bend/infrastructure/htmlTemplate"
	"fmt"
	"github.com/boltdb/bolt"
	"net/http"
)

func main() {
	db, err := bolt.Open("db/my.db", 0600, nil)
	requestRepository := boltDB.RequestRepository{DB: db}
	defer func() {
		_ = db.Close()
	}()

	http.Handle("/readme/", application.ReadmeService{})

	http.Handle("/dashboard/", application.DashboardService{
		RequestRepository:  requestRepository,
		DashboardPresenter: htmlTemplate.DashBoardPresenter{},
	})

	http.Handle("/sendRequest/", application.SendRequestService{
		RequestRepository: requestRepository,
		TransportService:  _http.Transport{},
	})

	http.Handle("/static/", http.FileServer(http.Dir("")))
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/favicon.ico")
	})

	http.Handle("/", application.RequestService{
		RequestRepository: requestRepository,
	})

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
