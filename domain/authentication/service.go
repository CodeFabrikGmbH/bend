package authentication

import (
	"context"
	"net/http"
)

type Service interface {
	Authenticate(w http.ResponseWriter, r *http.Request) (*User, error)
	Login(ctx context.Context, w http.ResponseWriter, username, password string) (*User, error)
	Logout(ctx context.Context, w http.ResponseWriter)
}
