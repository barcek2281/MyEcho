package storage

import "github.com/barcek2281/MyEcho/internal/app/model"


type BarcodeRepository struct {
	store *Storage
}


func (r *BarcodeRepository) Create(b *model.Barcode) error {
	err := r.store.db.QueryRow("INSERT INTO barcode (user_id, barcode) VALUES($1, $2) RETURNING id", b.User_id, b.Barcode).Scan(&b.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *BarcodeRepository) FindByUserId(user_id int) (*model.Barcode, error){
	b := &model.Barcode{}
	if err := r.store.db.QueryRow("SELECT id, user_id, barcode FROM barcode WHERE user_id = $1", user_id).Scan(&b.Id, &b.User_id, &b.Barcode); err != nil {
		return nil, err
	}
	return b, nil
}