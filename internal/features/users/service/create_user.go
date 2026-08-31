package users_service

import (
	"context"
	"fmt"

	"github.com/mercuryqa/todo-app/internal/core/domain"
)

// CreateUser создаёт нового пользователя: формирует доменный объект,
// валидирует его инварианты и сохраняет через репозиторий.
func (s *UsersService) CreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {

	if err := user.Validate(); err != nil {
		return domain.User{}, fmt.Errorf("validate user domain: %w", err)
	}

	user, err := s.usersRepository.CreateUser(ctx, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("save user in repository: %w", err)
	}

	return user, nil
}
