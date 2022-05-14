package repository

import (
	"database/sql"
	"errors"
	"math/rand"

	"github.com/google/uuid"

	"github.com/StenvL/async-architecture-course/tracker/app/model"
)

func (r Repository) GetUserTasks(userID int) (tasks []model.Task, err error) {
	query := "select id, title, status, description, assignee from tasks where assignee = $1"
	err = r.db.Select(&tasks, query, userID)

	return tasks, err
}

func (r Repository) GetTasksToShuffle() (tasks []model.Task, err error) {
	err = r.db.Select(&tasks, "select id, public_id, title, status, description, assignee from tasks where status = 'new'")
	return tasks, err
}

func (r Repository) CreateTask(task *model.Task) error {
	query := `
		insert into tasks (public_id, title, description, assignee)
		values ($1, $2, $3, (select id from users order by random() limit 1))
		returning id`

	var taskID int
	if err := r.db.QueryRow(query, task.PublicID, task.Title, task.Description).Scan(&taskID); err != nil {
		return err
	}

	return r.db.Get(
		task, "select id, public_id, title, status, description, created, assignee from tasks where id = $1", taskID,
	)
}

func (r Repository) MarkTaskAsResolved(userID, taskID int) (bool, uuid.UUID, error) {
	var publicID uuid.UUID
	query := "select public_id from tasks where id = $1 and assignee = $2;"
	if err := r.db.Get(&publicID, query, taskID, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, uuid.Nil, nil
		}
		return false, uuid.Nil, err
	}

	_, err := r.db.Exec("update tasks set status = 'resolved' where id = $1 and assignee = $2;", taskID, userID)
	return true, publicID, err
}

func (r Repository) ShuffleTasks(users []model.User, tasks []model.Task) (map[uuid.UUID]int, error) {
	res := make(map[uuid.UUID]int)

	tx, _ := r.db.Begin()
	for _, task := range tasks {
		userIdx := rand.Intn(len(users))
		res[task.PublicID] = users[userIdx].ID

		if _, err := tx.Exec("update tasks set assignee = $1 where id = $2", users[userIdx].ID, task.ID); err != nil {
			_ = tx.Rollback()
			return nil, err
		}
	}

	_ = tx.Commit()

	return res, nil
}
