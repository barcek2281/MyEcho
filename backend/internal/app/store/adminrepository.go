package storage

import "github.com/barcek2281/MyEcho/internal/app/model"

type AdminRepository struct {
	storage *Storage
}

func (r *AdminRepository) Create(a *model.Admin) error {
	err := a.BeforeCreate()
	if err != nil {
		return err
	}
	if err := r.storage.db.QueryRow("INSERT INTO admins (email, name, password) VALUES ($1, $2, $3) RETURNING id",
		a.Email, a.Name, a.Password,
	).Scan(&a.ID); err != nil {
		return err
	}
	return nil
}

func (r *AdminRepository) FindById(id int) (*model.Admin, error) {
	a := &model.Admin{}
	if err := r.storage.db.QueryRow("SELECT id, email, name, password FROM admins WHERE id = $1", id).Scan(
		&a.ID, &a.Email, &a.Name, &a.Password); err != nil {
		return nil, err
	}
	return a, nil
}

func (r *AdminRepository) FindByEmail(email string) (*model.Admin, error) {
	a := &model.Admin{}
	if err := r.storage.db.QueryRow("SELECT id, email, name, password FROM admins WHERE email = $1", email).Scan(
		&a.ID, &a.Email, &a.Name, &a.Password); err != nil {
		return nil, err
	}
	return a, nil
}
