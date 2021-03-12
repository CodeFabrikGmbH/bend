package environment

import (
	"code-fabrik.com/bend/domain/authentication"
	"code-fabrik.com/bend/domain/config"
	"code-fabrik.com/bend/domain/dashboardpage"
	"code-fabrik.com/bend/domain/loginpage"
	"code-fabrik.com/bend/domain/request"
)

type Environment struct {
	LoginPage         loginpage.Page
	DashboardPage     dashboardpage.Page
	RequestRepository request.Repository
	ConfigRepository  config.Repository
	Transport         request.Transport
	Authentication    authentication.Service
}
