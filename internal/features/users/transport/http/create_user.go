package users_transport_http

import (
	"net/http"

	"github.com/mercuryqa/todo-app/internal/core/domain"
	core_logger "github.com/mercuryqa/todo-app/internal/core/logger"
	core_http_response "github.com/mercuryqa/todo-app/internal/core/transport/http/http_response"
	core_http_request "github.com/mercuryqa/todo-app/internal/core/transport/http/request"
)

type CreateUserRequest struct {
	FullName    string  `json:"full_name" validate:"required,min=3,max=100"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15,startswith=+"`
}

type CreateUserResponse UserDTOResponse

// CreateUser - создает пользователя
func (h *UsersHTTPHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke CreateUser handler")
	// ...
	var request CreateUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	userDomain := domainFromDTO(request)
	userDomain, err := h.usersService.CreateUser(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create user")
		return
	}

	//if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
	//	fmt.Println("произошла ошибка")
	//}

	response := CreateUserResponse(userDTOFromDomain(userDomain))
	responseHandler.JSONResponse(response, http.StatusCreated)

	//w.WriteHeader(http.StatusOK)
}

func domainFromDTO(dto CreateUserRequest) domain.User {
	return domain.NewUserUnitialized(dto.FullName, dto.PhoneNumber)
}
