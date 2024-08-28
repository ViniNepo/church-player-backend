package sqlite

import (
	"database/sql"
	"errors"
	"log"

	"github.com/ViniNepo/church-player-backend/internal/domain"
)

type SectionDB struct {
	DB *sql.DB
}

func NewSectionDB(db *sql.DB) *SectionDB {
	return &SectionDB{
		DB: db,
	}
}

func (s *SectionDB) Create(section *domain.Section) error {
	if section == nil {
		return errors.New("section is nil")
	}

	var maxOrder *int
	err := s.DB.QueryRow("SELECT MAX(section_order) FROM sections WHERE worship_id=?", section.WorshipID).Scan(&maxOrder)
	if err != nil {
		log.Fatal(err)
	}

	var newOrder int
	if maxOrder != nil {
		newOrder = *maxOrder + 1
	} else {
		newOrder = 1
	}

	stmt, err := s.DB.Prepare("INSERT INTO sections (id, section_name, section_order, worship_id) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(section.ID, section.SectionName, newOrder, section.WorshipID)
	if err != nil {
		return err
	}

	return nil
}

func (s *SectionDB) FindAllByWorshipID(id string) (*[]domain.Section, error) {
	var sections []domain.Section
	rows, err := s.DB.Query("SELECT * FROM sections WHERE worship_id=? ORDER BY section_order ASC", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var section domain.Section
		err = rows.Scan(&section.ID, &section.SectionName, &section.SectionOrder, &section.WorshipID)
		if err != nil {
			return nil, err
		}

		sections = append(sections, section)
	}
	err = rows.Err()

	if err != nil {
		return nil, err
	}

	if len(sections) == 0 {
		return &[]domain.Section{}, nil
	}

	return &sections, nil
}

func (s *SectionDB) Update(section *domain.Section) error {
	if section == nil {
		return errors.New("section is nil")
	}

	stmt, err := s.DB.Prepare("UPDATE sections SET section_name=? WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(section.SectionName, section.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *SectionDB) Delete(id string) error {
	if id == "" {
		return errors.New("id is empty")
	}

	stmt, err := s.DB.Prepare("DELETE FROM sections WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	stmt, err = s.DB.Prepare("DELETE FROM moments WHERE section_id = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}
