package environment

import (
	"code-fabrik.com/bend/domain/config"
	"code-fabrik.com/bend/domain/request"
)

type Environment struct {
	RequestRepository request.Repository
	ConfigRepository  config.Repository
	Transport         request.Transport
}
