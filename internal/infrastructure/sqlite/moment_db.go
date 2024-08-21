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

	stmt, err := m.DB.Prepare("INSERT INTO moments (id, title, section_name, section_order, song_order, song_id, worship_id) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(moment.ID, moment.Title, moment.SectionName, moment.SectionOrder, moment.SongOrder, moment.SongID, moment.WorshipID)
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

	err = stmt.QueryRow(id).Scan(&moment.ID, &moment.Title, &moment.SectionName, &moment.SectionOrder, &moment.SongOrder, &moment.SongID, &moment.WorshipID)
	if err != nil {
		return nil, err
	}

	return &moment, nil
}

func (m *MomentDB) FindAllByWorshipID(id string) (*[]domain.Moment, error) {
	if id == "" {
		return nil, errors.New("worship id is empty")
	}

	var moments []domain.Moment
	rows, err := m.DB.Query("SELECT * FROM moments WHERE worship_id = ? ORDER BY song_order", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var moment domain.Moment
		err = rows.Scan(&moment.ID, &moment.Title, &moment.SectionName, &moment.SectionOrder, &moment.SongOrder, &moment.SongID, &moment.WorshipID)
		if err != nil {
			return nil, err
		}

		moments = append(moments, moment)
	}
	err = rows.Err()

	if err != nil {
		return nil, err
	}

	if len(moments) == 0 {
		return &[]domain.Moment{}, nil
	}

	return &moments, nil
}

func (m *MomentDB) Update(moment *domain.Moment) error {
	if moment == nil {
		return errors.New("album is nil")
	}

	stmt, err := m.DB.Prepare("UPDATE moments SET title=?, section_name=?, song_id=? WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(moment.Title, moment.SectionName, moment.SongID, moment.ID)
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

	stmt, err := tx.Prepare("UPDATE moments SET song_order = ? WHERE id = ?")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, update := range *momentOrderDTO {
		_, err := stmt.Exec(update.SongOrder, update.ID)
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
