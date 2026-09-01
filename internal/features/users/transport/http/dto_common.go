package users_transport_http

import (
	"github.com/mercuryqa/todo-app/internal/core/domain"
)

// UserDTOResponse — DTO для представления пользователя в API-ответе.
type UserDTOResponse struct {
	ID          int     `json:"id"           example:"1"`
	Version     int     `json:"version"      example:"3"`
	FullName    string  `json:"full_name"    example:"Ivan Ivanov"`
	PhoneNumber *string `json:"phone_number" example:"+79998887766"`
}

// userDTOFromDomain конвертирует доменный объект User в DTO для HTTP-ответа.
func userDTOFromDomain(user domain.User) UserDTOResponse {
	return UserDTOResponse{
		ID:          user.ID,
		Version:     user.Version,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	}
}

// usersDTOFromDomains конвертирует список доменных объектов в список DTO.
func usersDTOFromDomains(users []domain.User) []UserDTOResponse {
	usersDTO := make([]UserDTOResponse, len(users))

	for i, user := range users {
		usersDTO[i] = userDTOFromDomain(user)
	}

	return usersDTO
}
