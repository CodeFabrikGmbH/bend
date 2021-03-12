package application

import (
	"code-fabrik.com/bend/domain/environment"
)

type DeletionService struct {
	Env environment.Environment
}

func (ds DeletionService) Delete(path string, requestId *string) error {
	if requestId == nil {
		return ds.Env.RequestRepository.DeletePath(path)
	} else {
		return ds.Env.RequestRepository.DeleteRequestForPath(path, *requestId)
	}
}
