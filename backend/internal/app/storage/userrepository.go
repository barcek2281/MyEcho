package storage

import "github.com/barcek2281/MyEcho/internal/app/model"

// UserRepo
type UserRepository struct {
	storage *Storage
}

func (r *UserRepository) Create(u *model.User) (*model.User, error) {
	if err := r.storage.db.QueryRow("INSERT INTO users (email, login, password) VALUES ($1, $2, $3) RETURNING id)",
		u.Email, u.Login, u.Password,
	).Scan(&u.ID); err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	if err := r.storage.db.QueryRow("SELECT id, email, login, password FROM users WHERE email = $1",
		email).Scan(&u.ID, &u.Email, &u.Login, &u.Password); err != nil {
		return nil, err
	}
	return nil, nil
}
