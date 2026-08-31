package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	core_errors "github.com/mercuryqa/todo-app/internal/core/errors"
)

// requestValidator — глобальный экземпляр валидатора (go-playground/validator).
// Используется для валидации struct-тегов вида `validate:"required,min=1,max=100"`.
// Создаётся один раз на пакет (потокобезопасен).
var requestValidator = validator.New()

// DecodeAndValidateRequest десериализует тело запроса из JSON в dest
// и затем проверяет валидность данных.
//
// Порядок валидации:
//  1. Если dest реализует validatable — вызывается его Validate()
//  2. Иначе — валидируются struct-теги через go-playground/validator
func DecodeAndValidateRequest(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf(
			"decode json: %v: %w",
			err,
			core_errors.ErrInvalidArgument)
	}

	if err := requestValidator.Struct(dest); err != nil {
		return fmt.Errorf(
			"request validation: %v: %w",
			err,
			core_errors.ErrInvalidArgument)
	}

	return nil
}
