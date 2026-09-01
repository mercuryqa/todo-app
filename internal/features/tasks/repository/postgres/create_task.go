package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/mercuryqa/todo-app/internal/core/domain"
	core_errors "github.com/mercuryqa/todo-app/internal/core/errors"
	core_postgres_pool "github.com/mercuryqa/todo-app/internal/core/repository/postgres/pool"
)

// SaveTask вставляет новую задачу в БД и возвращает сохранённую версию (с данными из БД).
// RETURNING позволяет получить итоговое состояние записи одним запросом,
// без дополнительного SELECT (т.н. «insert and return» паттерн).
//
// Если author_user_id не существует в таблице users — PostgreSQL вернёт
// ошибку внешнего ключа (код 23503), которую адаптер pgx преобразует
// в ErrViolatesForeignKey → мы транслируем в ErrNotFound (автор не найден).
func (r *TasksRepository) SaveTask(
	ctx context.Context,
	task domain.Task,
) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO todoapp.tasks (title, description, completed, created_at, completed_at, author_user_id)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, version, title, description, completed, created_at, completed_at, author_user_id;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Completed,
		task.CreatedAt,
		task.CompletedAt,
		task.AuthorUserID,
	)

	var taskModel TaskModel
	if err := taskModel.Scan(row); err != nil {
		if errors.Is(err, core_postgres_pool.ErrViolatesForeignKey) {
			return domain.Task{}, fmt.Errorf(
				"%v: user with id='%s': %w",
				err,
				task.AuthorUserID,
				core_errors.ErrNotFound,
			)
		}

		return domain.Task{}, fmt.Errorf("scan error: %w", err)
	}

	taskDomain := modelToDomain(taskModel)

	return taskDomain, nil
}
