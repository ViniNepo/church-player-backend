package sqlite

import (
	"database/sql"
	"errors"
	"mime/multipart"

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

func (s *SongDB) Create(song *domain.Song, fileHeader *multipart.FileHeader) error {
	if song == nil {
		return errors.New("song is nil")
	}

	var maxOrder *int
	err := s.DB.QueryRow("SELECT MAX(song_order) FROM songs WHERE album_id=?", song.AlbumID).Scan(&maxOrder)
	if err != nil {
		return err
	}

	var newOrder int
	if maxOrder != nil {
		newOrder = *maxOrder + 1
	} else {
		newOrder = 1
	}

	stmt, err := s.DB.Prepare("INSERT INTO songs (id, title, song_number, file, song_order, album_id) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(song.ID, song.Title, song.Number, song.File, newOrder, song.AlbumID)
	if err != nil {
		return err
	}

	return nil
}

func (s *SongDB) FindAll() (*[]dto.SongWithAlbumDTO, error) {
	var songs []dto.SongWithAlbumDTO
	rows, err := s.DB.Query("SELECT s.id, s.title, s.song_number, s.file, s.song_order, a.id, a.title, a.image_url FROM songs as s INNER JOIN albums as a ON s.album_id = a.id ORDER BY s.title ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var song dto.SongWithAlbumDTO
		err := rows.Scan(&song.ID, &song.Title, &song.SongNumber, &song.File, &song.Order, &song.AlbumID, &song.AlbumName, &song.AlbumImageURL)
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
		return &[]dto.SongWithAlbumDTO{}, nil
	}

	return &songs, nil
}

func (s *SongDB) FindAllByAlbumID(id string) (*[]domain.Song, error) {
	var songs []domain.Song
	rows, err := s.DB.Query("SELECT * FROM songs WHERE album_id=? ORDER BY song_order", id)
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

	stmt, err := s.DB.Prepare("UPDATE songs SET title=?, file=?, song_number=?, song_order=?, album_id=? WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(song.Title, song.File, song.Number, song.Order, song.AlbumID, song.ID)
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

	stmt, err := tx.Prepare("UPDATE songs SET song_order = ? WHERE id = ?")
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
