package storage

import (
	"errors"
	"fmt"

	"github.com/barcek2281/MyEcho/internal/app/model"
)

var errEmailIsUsed = errors.New("email: email is used")

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
		return errEmailIsUsed
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

func (r *UserRepository) FindById(id int) (*model.User, error) {
	u := &model.User{}
	if err := r.storage.db.QueryRow("SELECT id, email, login, password FROM users WHERE id = $1", id).Scan(&u.ID,
		&u.Email, &u.Login, &u.Password); err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) ChangeLoginByEmail(newLogin, email string) error {
	// Выполнение запроса для обновления записи
	result, err := r.storage.db.Exec("UPDATE users SET login = $1 WHERE email = $2", newLogin, email)
	if err != nil {
		// Возвращаем ошибку, если запрос не удался
		return fmt.Errorf("failed to update login for email %s: %w", email, err)
	}

	// Проверяем, сколько строк было обновлено
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no user found with email %s", email)
	}

	return nil
}

func (r *UserRepository) DeleteByEmail(email string) error {
	result, err := r.storage.db.Exec("DELETE FROM users WHERE email = $1", email)
	if err != nil {
		// Возвращаем ошибку, если запрос не удался
		return fmt.Errorf("failed to delete user %s: %w", email, err)
	}

	// Проверяем, сколько строк было обновлено
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no user found with email %s", email)
	}

	return nil
}

func (r *UserRepository) GetAll(limit int) ([]*model.User, error) {
	rows, err := r.storage.db.Query("SELECT id, email, login, password FROM users ORDER BY id LIMIT $1", limit)
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

func (r *UserRepository) GetAllWithoutLimit() ([]*model.User, error) {
	rows, err := r.storage.db.Query("SELECT id, email, login, password FROM users ORDER BY id")
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

func (r *UserRepository) Activate(id int) error {
	result, err := r.storage.db.Exec("UPDATE users SET is_active = true WHERE id = $1", id)
	if err != nil {
		return err
	}
	// Проверяем, сколько строк было обновлено
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("not found user")
	}

	return nil
}

func (r *UserRepository) IsActive(id int) (*model.User, error) {
	u := &model.User{}
	if err := r.storage.db.QueryRow("SELECT id, email, login, password FROM users WHERE id = $1 AND is_active = true", id).Scan(&u.ID,
		&u.Email, &u.Login, &u.Password); err != nil {
		return nil, err
	}
	return u, nil
}
