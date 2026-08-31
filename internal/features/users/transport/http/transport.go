package users_transport_http

import (
	"context"
	"net/http"

	"github.com/mercuryqa/todo-app/internal/core/domain"
	core_http_middleware "github.com/mercuryqa/todo-app/internal/core/transport/http/middleware"
	"github.com/mercuryqa/todo-app/internal/core/transport/http/server"
)

type UsersHTTPHandler struct {
	usersService usersService
}

type usersService interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUsers(ctx context.Context, limit *int, offset *int) ([]domain.User, error)
	GetUser(ctx context.Context, id int) (domain.User, error)
	DeleteUser(ctx context.Context, id int) error
	PatchUser(ctx context.Context, id int, patch domain.UserPatch) (domain.User, error)
}

func NewUsersHTTPHandler(usersService usersService) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		usersService: usersService,
	}
}

func (h *UsersHTTPHandler) Routes() []server.Route {
	return []server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: h.CreateUser,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: h.GetUsers,
			// Пример точечного добавления middleware на V1
			Middleware: []core_http_middleware.Middleware{
				core_http_middleware.Dummy("get users middleware"),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/{id}",
			Handler: h.GetUser,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/users/{id}",
			Handler: h.DeleteUser,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/users/{id}",
			Handler: h.PatchUser,
		},
	}
}
