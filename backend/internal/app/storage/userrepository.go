package storage

import (
	"github.com/barcek2281/MyEcho/internal/app/model"
	"log"
)

// UserRepository
type UserRepository struct {
	storage *Storage
}

func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {

		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	if err := r.storage.db.QueryRow("INSERT INTO users (email, login, password) VALUES ($1, $2, $3) RETURNING id",
		u.Email, u.Login, u.Password,
	).Scan(&u.ID); err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	if err := r.storage.db.QueryRow("SELECT id, email, login, password FROM users WHERE email = $1",
		email).Scan(&u.ID, &u.Email, &u.Login, &u.Password); err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) ChangeLoginByEmail(newLogin, email, encPassword string) (*model.User, error) {
	u := &model.User{}
	if err := r.storage.db.QueryRow("UPDATE users SET login = $1 WHere email = $2 AND password = $3", newLogin, email, encPassword).Scan(&u.ID); err != nil {
		log.Fatalln("CANNOT FIND SUCH A USER", err)
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) DeleteByEmailAndPasswd(email, encPassword string) error {
	if err := r.storage.db.QueryRow("DELETE FROM users WHERE email = $1 AND password = $2", email, encPassword).Scan(); err != nil {
		log.Fatalln("CANNOT FIND SUCH A USER", err)
		return err
	}
	return nil
}

func (r *UserRepository) GetAll(limit int) ([]*model.User, error) {
	rows, err := r.storage.db.Query("SELECT id, email, login, password FROM users LIMIT $1", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*model.User
	for rows.Next() {
		u := &model.User{}
		if err := rows.Scan(&u.ID, &u.Email, &u.Login, &u.Password); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
