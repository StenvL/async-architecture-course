package repository

import (
	"github.com/StenvL/async-architecture-course/billing/app/model"
)

func (r Repository) GetUsers() (users []model.User, err error) {
	err = r.db.Select(&users, "select id, name, role from users")
	return users, err
}

func (r Repository) CreateUser(user model.User) error {
	query := `
		insert into users (id, name, email, role)
		values (:id, :name, :email, :role)`
	_, err := r.db.NamedExec(query, user)

	return err
}

func (r Repository) UpdateUser(user model.User) error {
	query := `
		update users set
			name = :name,
		    email = :email
		where id = :id`
	_, err := r.db.NamedExec(query, user)

	return err
}

func (r Repository) UpdateUserRole(id int, role string) error {
	_, err := r.db.Exec("update users set role = $1 where id = $2", role, id)
	return err
}

func (r Repository) DeleteUser(id int) error {
	_, err := r.db.Exec("delete from users where id = $1", id)
	return err
}
