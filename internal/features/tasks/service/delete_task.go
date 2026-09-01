package tasks_service

import (
	"context"
	"fmt"
)

// DeleteTask удаляет задачу по ID, делегируя запрос репозиторию.
func (s *TasksService) DeleteTask(
	ctx context.Context,
	id int,
) error {
	if err := s.tasksRepository.DeleteTask(ctx, id); err != nil {
		return fmt.Errorf("delete task from repository: %w", err)
	}

	return nil
}
