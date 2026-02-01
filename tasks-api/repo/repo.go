package repo

import (
	"context"
	"tasks-api/models"

	"github.com/jackc/pgx/v5"
)

type Repo struct {
	db *pgx.Conn
}

func NewRepo(db *pgx.Conn) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) Create(ctx context.Context, task models.CreateTask) error {
	_, err := r.db.Exec(
		ctx,
		"INSERT INTO tasks (name, description) VALUES ($1, $2)",
		task.Name, task.Description,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) RequestTasks(ctx context.Context) ([]models.TaskResponse, error) {
	var tasks []models.TaskResponse

	rows, err := r.db.Query(ctx, "SELECT id, name, description, created_at FROM tasks")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var task models.TaskResponse
		err := rows.Scan(&task.Id, &task.Name, &task.Description, &task.CreatedAt)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *Repo) RequestTask(ctx context.Context, taskId string) (models.TaskResponse, error) {
	var task models.TaskResponse

	row := r.db.QueryRow(
		ctx,
		"SELECT id, name, description, created_at FROM tasks WHERE id=$1",
		taskId,
	)
	err := row.Scan(&task.Id, &task.Name, &task.Description, &task.CreatedAt)
	if err != nil {
		return models.TaskResponse{}, err
	}

	return task, nil
}

func (r *Repo) Delete(ctx context.Context, taskId string) error {
	_, err := r.db.Exec(ctx, "DELETE FROM tasks WHERE id=$1", taskId)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) Change(ctx context.Context, task models.CreateTask, taskId string) error {
	_, err := r.db.Exec(
		ctx,
		"UPDATE tasks SET name = $1, description = $2 WHERE id = $3",
		task.Name, task.Description, taskId,
	)
	if err != nil {
		return err
	}

	return nil
}
