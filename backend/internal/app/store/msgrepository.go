package storage

import "github.com/barcek2281/MyEcho/internal/app/model"

type MsgRepository struct {
	store *Storage
}

func (m *MsgRepository) CreateMessage(msg *model.Messages) error {

	err := m.store.db.QueryRow("INSERT INTO messages (sender, receiver, message) VALUES($1, $2, $3) RETURING id", msg.Sender, msg.Receiver, msg.Message).Scan(&msg.Id)
	if err != nil {
		return err
	}
	return nil
}
