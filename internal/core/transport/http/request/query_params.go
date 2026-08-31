package core_http_utils

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/mercuryqa/todo-app/internal/core/errors"
)

// GetIntQueryParam - получает query параметры из URL
func GetIntQueryParam(r *http.Request, key string) (*int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}
	val, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%5' by key='%s' not a valid integer %v: %w",
			param,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return &val, nil

}
