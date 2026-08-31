// Package domain содержит основные доменные сущности приложения.
package domain

import (
	"fmt"
	"regexp"

	core_errors "github.com/mercuryqa/todo-app/internal/core/errors"
)

// User представляет пользователя системы.
type User struct {
	ID      int
	Version int

	FullName    string
	PhoneNumber *string
}

// NewUser — конструктор для восстановления пользователя по имеющему набору данных
func NewUser(
	id int,
	version int,
	fullName string,
	phoneNumber *string,
) User {
	return User{
		ID:          id,
		Version:     version,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}

// NewUserUnitialized — конструктор для восстановления пользователя по имеющему набору данных
func NewUserUnitialized(
	fullName string,
	phoneNumber *string,
) User {
	return NewUser(
		UnitializedID,
		UnitializecVersion,
		fullName,
		phoneNumber,
	)
}

// Validate проверяет инварианты пользователя.
// Формат телефона: начинается с «+», далее только цифры, длина 10–15 символов.
// Пример: +79001234567
func (u *User) Validate() error {
	fullNameLen := len([]rune(u.FullName))
	if fullNameLen < 3 || fullNameLen > 100 {
		return fmt.Errorf(
			"invalid `FullName` len: %d: %w",
			fullNameLen,
			core_errors.ErrInvalidArgument,
		)
	}

	if u.PhoneNumber != nil {
		phoneNumberLen := len([]rune(*u.PhoneNumber))
		if phoneNumberLen < 10 || phoneNumberLen > 15 {
			return fmt.Errorf(
				"invalid `PhoneNumber` len: %d: %w",
				phoneNumberLen,
				core_errors.ErrInvalidArgument,
			)
		}

		// regexp.MustCompile паникует при невалидном паттерне — это допустимо,
		// так как паттерн — константа, известная на этапе компиляции.
		// В продакшн-коде регулярное выражение лучше вынести в переменную пакета.
		re := regexp.MustCompile(`^\+[0-9]+$`)

		if !re.MatchString(*u.PhoneNumber) {
			return fmt.Errorf(
				"invalid `PhoneNumber` format: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	return nil
}

// UserPatch содержит изменения для частичного обновления пользователя (PATCH).
// Каждое поле обёрнуто в Nullable, чтобы различать «не передано» и «передано null».
// Подробнее о Nullable: см. internal/core/domain/nullable.go.
type UserPatch struct {
	FullName    Nullable[string]
	PhoneNumber Nullable[string]
}

// NewUserPatch — конструктор UserPatch.
func NewUserPatch(
	fullName Nullable[string],
	phoneNumber Nullable[string],
) UserPatch {
	return UserPatch{
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}

// Validate проверяет корректность патча до его применения.
// FullName является обязательным полем и не может быть обнулён.
func (p *UserPatch) Validate() error {
	if p.FullName.Set && p.FullName.Value == nil {
		return fmt.Errorf(
			"`FullName` can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}

// ApplyPatch применяет изменения к пользователю.
// Используется та же техника «копия → изменение → валидация → замена», что и в Task.ApplyPatch.
func (u *User) ApplyPatch(patch UserPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate user patch: %w", err)
	}

	tmp := *u

	// Проводит валидацию и если что-то не так откатываем возрвращаем юзера из временной переменной tmp
	if patch.FullName.Set {
		tmp.FullName = *patch.FullName.Value
	}

	if patch.PhoneNumber.Set {
		tmp.PhoneNumber = patch.PhoneNumber.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched user: %w", err)
	}

	*u = tmp

	return nil
}
