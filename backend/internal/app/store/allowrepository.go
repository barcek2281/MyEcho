package storage

import (
	"log"

	"github.com/barcek2281/MyEcho/internal/app/model"
)

type AllowRepository struct {
	store *Storage
}

func (a *AllowRepository) GetAllow() (map[string]map[string]bool, error) {
	res := make(map[string]map[string]bool)
	query := `SELECT * FROM allow`
	rows, err := a.store.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var allow model.Allow
		if err := rows.Scan(&allow.Id, &allow.EmailFirst, &allow.EmailSecond); err != nil {
			log.Fatal("error:", err)
			continue
		}
		res[allow.EmailFirst][allow.EmailSecond] = true
		res[allow.EmailSecond][allow.EmailFirst] = true
	}
	return res, nil
}

func (a *AllowRepository) RemoveAllow(allow model.Allow) error {
	_, err := a.store.db.Query("DELETE FROM allow WHERE (emailFirst = $1 AND emailSecond = $2) OR (emailSecond = $3 OR emailFirst = $4)", allow.EmailFirst, allow.EmailSecond, allow.EmailFirst, allow.EmailSecond)
	if err != nil {
		return err
	}
	return nil
}

func (a *AllowRepository) AddAllow(allow model.Allow) error {
	_, err := a.store.db.Query("INSERT INTO allow (emailFirst, emailSecond) VALUES ($1, $2)", allow.EmailFirst, allow.EmailSecond)
	if err != nil {
		return err
	}
	return nil
}
