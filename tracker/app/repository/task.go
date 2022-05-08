package repository

import (
	"math/rand"

	"github.com/StenvL/async-architecture-course/app/model"
)

func (r Repository) GetUserTasks(userID int) (tasks []model.Task, err error) {
	query := "select id, title, status, description, assignee from tasks where assignee = $1"
	err = r.db.Select(&tasks, query, userID)

	return tasks, err
}

func (r Repository) GetTaskIDsToShuffle() (ids []int, err error) {
	err = r.db.Select(&ids, "select id from tasks where status = 'new'")
	return ids, err
}

func (r Repository) CreateTask(task model.Task) error {
	query := `
		insert into tasks (title, description, assignee)
		values (:title, :description, (select id from users order by random() limit 1))`
	_, err := r.db.NamedExec(query, task)

	return err
}

func (r Repository) MarkTaskAsResolved(userID, taskID int) error {
	_, err := r.db.Exec("update tasks set status = 'resolved' where id = $1 and assignee = $2;", taskID, userID)
	return err
}

func (r Repository) ShuffleTasks(users []model.User, taskIDs []int) (map[int]model.User, error) {
	res := make(map[int]model.User)

	tx, _ := r.db.Begin()
	for _, taskID := range taskIDs {
		userIdx := rand.Intn(len(users))
		res[taskID] = users[userIdx]

		if _, err := tx.Exec("update tasks set assignee = $1 where id = $2", users[userIdx].ID, taskID); err != nil {
			_ = tx.Rollback()
			return nil, err
		}
	}

	_ = tx.Commit()

	return res, nil
}
