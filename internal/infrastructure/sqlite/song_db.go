package sqlite

import (
	"database/sql"
	"errors"

	"github.com/ViniNepo/church-player-backend/internal/domain"
	"github.com/ViniNepo/church-player-backend/internal/domain/dto"
)

type SongDB struct {
	DB *sql.DB
}

func NewSongDB(db *sql.DB) *SongDB {
	return &SongDB{
		DB: db,
	}
}

func (s *SongDB) Create(song *domain.Song) error {
	if song == nil {
		return errors.New("song is nil")
	}

	stmt, err := s.DB.Prepare("INSERT INTO songs (id, title, number, file, order_song, album_id) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(song.ID, song.Title, song.Number, song.File, song.Order, song.AlbumID)
	if err != nil {
		return err
	}

	return nil
}

func (s *SongDB) FindAll() (*[]domain.Song, error) {
	var songs []domain.Song
	rows, err := s.DB.Query("SELECT * FROM songs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var song domain.Song
		err := rows.Scan(&song.ID, &song.Title, &song.Number, &song.File, &song.Order, &song.AlbumID)
		if err != nil {
			return nil, err
		}

		songs = append(songs, song)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	if len(songs) == 0 {
		return &[]domain.Song{}, nil
	}

	return &songs, nil
}

func (s *SongDB) FindAllByAlbumID(id string) (*[]domain.Song, error) {
	var songs []domain.Song
	rows, err := s.DB.Query("SELECT * FROM songs WHERE album_id=?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var song domain.Song
		err := rows.Scan(&song.ID, &song.Title, &song.Number, &song.File, &song.Order, &song.AlbumID)
		if err != nil {
			return nil, err
		}

		songs = append(songs, song)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &songs, nil
}

func (s *SongDB) FindByID(id string) (*domain.Song, error) {
	var song domain.Song
	stmt, err := s.DB.Prepare("SELECT * FROM songs WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&song.ID, &song.Title, &song.Number, &song.File, &song.Order, &song.AlbumID)
	if err != nil {
		return nil, err
	}

	return &song, nil
}

func (s *SongDB) Update(song *domain.Song) error {
	if song == nil {
		return errors.New("song is nil")
	}

	stmt, err := s.DB.Prepare("UPDATE songs SET title=?, number=?, file=?, order_song=?, album_id=? WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(song.Title, song.Number, song.File, song.Order, song.AlbumID, song.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *SongDB) UpdateOrder(songUpdateOrderDTO *[]dto.SongUpdateOrderDTO) error {
	if songUpdateOrderDTO == nil {
		return errors.New("songUpdateOrderDTO is nil")
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("UPDATE songs SET order_song = ? WHERE id = ?")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, update := range *songUpdateOrderDTO {
		_, err := stmt.Exec(update.Order, update.ID)
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

func (s *SongDB) Delete(id string) error {
	if id == "" {
		return errors.New("song is nil")
	}

	stmt, err := s.DB.Prepare("DELETE FROM songs WHERE id=?")
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
