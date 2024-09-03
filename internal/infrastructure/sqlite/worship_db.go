package sqlite

import (
	"database/sql"
	"errors"
	"github.com/ViniNepo/church-player-backend/internal/domain"
)

type WorshipDB struct {
	DB *sql.DB
}

func NewWorshipDB(db *sql.DB) *WorshipDB {
	return &WorshipDB{
		DB: db,
	}
}

func (w *WorshipDB) Create(worship *domain.Worship) error {
	if worship == nil {
		return errors.New("worship is nil")
	}

	var maxOrder *int
	err := w.DB.QueryRow("SELECT MAX(worship_order) FROM worships").Scan(&maxOrder)
	if err != nil {
		return err
	}

	var newOrder int
	if maxOrder != nil {
		newOrder = *maxOrder + 1
	} else {
		newOrder = 1
	}

	stmt, err := w.DB.Prepare("INSERT INTO worships (id, title, image_url, worship_order) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(worship.ID, worship.Title, worship.ImageUrl, newOrder)
	if err != nil {
		return err
	}

	return nil
}

func (w *WorshipDB) FindByID(id string) (*domain.Worship, error) {
	if id == "" {
		return nil, errors.New("id is empty")
	}

	var worship domain.Worship
	stmt, err := w.DB.Prepare("SELECT * FROM worships WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&worship.ID, &worship.Title, &worship.ImageUrl, &worship.WorshipOrder)
	if err != nil {
		return nil, err
	}

	return &worship, nil
}

func (w *WorshipDB) FindAll() (*[]domain.Worship, error) {
	var worships []domain.Worship
	rows, err := w.DB.Query("SELECT * FROM worships ORDER BY worship_order ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var worship domain.Worship
		err = rows.Scan(&worship.ID, &worship.Title, &worship.ImageUrl, &worship.WorshipOrder)
		if err != nil {
			return nil, err
		}

		worships = append(worships, worship)
	}
	err = rows.Err()

	if err != nil {
		return nil, err
	}

	if len(worships) == 0 {
		return &[]domain.Worship{}, nil
	}

	return &worships, nil
}

func (w *WorshipDB) Update(worship *domain.Worship) error {
	if worship == nil {
		return errors.New("worship is nil")
	}

	stmt, err := w.DB.Prepare("UPDATE worships SET title=?, image_url=? WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(worship.Title, worship.ImageUrl, worship.ID)
	if err != nil {
		return err
	}

	return nil
}

func (w *WorshipDB) Delete(id string) error {
	if id == "" {
		return errors.New("id is empty")
	}

	stmt, err := w.DB.Prepare("DELETE FROM worships WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}
