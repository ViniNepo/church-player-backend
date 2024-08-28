package sqlite

import (
	"database/sql"
	"errors"
	"log"

	"github.com/ViniNepo/church-player-backend/internal/domain"
)

type AlbumDB struct {
	DB *sql.DB
}

func NewAlbumDB(db *sql.DB) *AlbumDB {
	return &AlbumDB{
		DB: db,
	}
}

func (a *AlbumDB) Create(album *domain.Album) error {
	if album == nil {
		return errors.New("album is nil")
	}

	var maxOrder *int
	err := a.DB.QueryRow("SELECT MAX(album_order) FROM albums").Scan(&maxOrder)
	if err != nil {
		log.Fatal(err)
	}

	var newOrder int
	if maxOrder != nil {
		newOrder = *maxOrder + 1
	} else {
		newOrder = 1
	}

	stmt, err := a.DB.Prepare("INSERT INTO albums (id, title, image_url, album_order) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(album.ID, album.Title, album.ImageUrl, newOrder)
	if err != nil {
		return err
	}

	return nil
}

func (a *AlbumDB) FindByID(id string) (*domain.Album, error) {
	if id == "" {
		return nil, errors.New("id is empty")
	}

	var album domain.Album
	stmt, err := a.DB.Prepare("SELECT * FROM albums WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&album.ID, &album.Title, &album.ImageUrl, &album.AlbumOrder)
	if err != nil {
		return nil, err
	}

	return &album, nil
}

func (a *AlbumDB) FindAll() (*[]domain.Album, error) {
	var albums []domain.Album
	rows, err := a.DB.Query("SELECT * FROM albums ORDER BY album_order ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var album domain.Album
		err = rows.Scan(&album.ID, &album.Title, &album.ImageUrl, &album.AlbumOrder)
		if err != nil {
			return nil, err
		}

		albums = append(albums, album)
	}
	err = rows.Err()

	if err != nil {
		return nil, err
	}

	if len(albums) == 0 {
		return &[]domain.Album{}, nil
	}

	return &albums, nil
}

func (a *AlbumDB) Update(album *domain.Album) error {
	if album == nil {
		return errors.New("album is nil")
	}

	stmt, err := a.DB.Prepare("UPDATE albums SET title=?, image_url=? WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(album.Title, album.ImageUrl, album.ID)
	if err != nil {
		return err
	}

	return nil
}

func (a *AlbumDB) Delete(id string) error {
	if id == "" {
		return errors.New("id is empty")
	}

	stmt, err := a.DB.Prepare("DELETE FROM albums WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	stmt, err = a.DB.Prepare("DELETE FROM songs WHERE album_id = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}
