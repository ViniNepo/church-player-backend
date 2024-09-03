package sqlite

import (
	"database/sql"
	"errors"
	"github.com/ViniNepo/church-player-backend/internal/domain"
	"github.com/ViniNepo/church-player-backend/internal/domain/dto"
)

type MomentDB struct {
	DB *sql.DB
}

func NewMomentDB(db *sql.DB) *MomentDB {
	return &MomentDB{
		DB: db,
	}
}

func (m *MomentDB) Create(moment *domain.Moment) error {
	if moment == nil {
		return errors.New("moment is nil")
	}

	var maxOrder *int
	err := m.DB.QueryRow("SELECT MAX(song_order) FROM moments WHERE section_id=?", moment.SectionID).Scan(&maxOrder)
	if err != nil {
		return err
	}

	var newOrder int
	if maxOrder != nil {
		newOrder = *maxOrder + 1
	} else {
		newOrder = 1
	}

	stmt, err := m.DB.Prepare("INSERT INTO moments (id, title, song_order, song_id, section_id) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(moment.ID, moment.Title, newOrder, moment.SongID, moment.SectionID)
	if err != nil {
		return err
	}
	return nil
}

func (m *MomentDB) FindByID(id string) (*domain.Moment, error) {
	if id == "" {
		return nil, errors.New("moment id is empty")
	}

	var moment domain.Moment
	stmt, err := m.DB.Prepare("SELECT * FROM moments WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&moment.ID, &moment.Title, &moment.SongOrder, &moment.SongID, &moment.SectionID)
	if err != nil {
		return nil, err
	}

	return &moment, nil
}

func (m *MomentDB) FindAllBySectionID(id string) (*[]domain.Moment, error) {
	var moments []domain.Moment
	rows, err := m.DB.Query("SELECT * FROM moments WHERE section_id=? ORDER BY song_order", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var moment domain.Moment
		err := rows.Scan(&moment.ID, &moment.Title, &moment.SongOrder, &moment.SongID, &moment.SectionID)
		if err != nil {
			return nil, err
		}

		moments = append(moments, moment)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &moments, nil
}

func (m *MomentDB) Update(moment *domain.Moment) error {
	if moment == nil {
		return errors.New("album is nil")
	}

	stmt, err := m.DB.Prepare("UPDATE moments SET title=?, song_id=? WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(moment.Title, moment.SongID, moment.ID)
	if err != nil {
		return err
	}

	return nil
}

func (m *MomentDB) UpdateOrder(momentOrderDTO *[]dto.MomentOrderDTO) error {
	if momentOrderDTO == nil {
		return errors.New("momentOrderDTO is nil")
	}

	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("UPDATE moments SET song_order = ?, section_id=? WHERE id = ?")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, update := range *momentOrderDTO {
		_, err := stmt.Exec(update.SongOrder, update.SectionID, update.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (m *MomentDB) Delete(id string) error {
	if id == "" {
		return errors.New("moment is nil")
	}

	stmt, err := m.DB.Prepare("DELETE FROM moments WHERE id=?")
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
