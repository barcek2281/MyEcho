package storage

import (
	"log"

	"github.com/barcek2281/MyEcho/internal/app/model"
)

type MsgRepository struct {
	store *Storage
}

func (m *MsgRepository) CreateMessage(msg *model.Messages) error {
	err := m.store.db.QueryRow("INSERT INTO messages (sender, receiver, msg) VALUES($1, $2, $3) RETURNING id", msg.Sender, msg.Receiver, msg.Message).Scan(&msg.Id)
	if err != nil {
		return err
	}
	return nil
}

func (m *MsgRepository) GetMsg(receiver, sender string, limit int) ([]*model.Messages, error) {
	var res []*model.Messages

	query := `SELECT * FROM (
			SELECT id, sender, receiver, msg, date 
			FROM messages
			WHERE (sender = $1 AND receiver = $2) OR (sender = $2 AND receiver = $1)
			ORDER BY date DESC
			LIMIT $3
		) AS subquery ORDER BY date ASC;`

	rows, err := m.store.db.Query(query, receiver, sender, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var msg model.Messages
		if err := rows.Scan(&msg.Id, &msg.Sender, &msg.Receiver, &msg.Message, &msg.Date); err != nil {
			log.Fatal("error to get data from database messages", err)
			continue
		}
		res = append(res, &msg)

	}

	return res, nil
}
