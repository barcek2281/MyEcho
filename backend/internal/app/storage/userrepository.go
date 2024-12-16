package storage

import (
	"fmt"
	"github.com/barcek2281/MyEcho/internal/app/model"
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
