package sqlite_test

import (
	"database/sql"
	"testing"

	"github.com/ViniNepo/church-player-backend/internal/domain"
	"github.com/ViniNepo/church-player-backend/internal/domain/dto"
	"github.com/ViniNepo/church-player-backend/internal/infrastructure/sqlite"
	"github.com/stretchr/testify/suite"
)

type SongDBTestSuite struct {
	suite.Suite
	db      *sql.DB
	songDB  *sqlite.SongDB
	albumDB *sqlite.AlbumDB
}

func (s *SongDBTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.db = db

	db.Exec("CREATE table albums (id varchar(255), title varchar(255), image_url varchar(255))")
	s.albumDB = sqlite.NewAlbumDB(db)
	db.Exec("CREATE table songs (id varchar(255), title varchar(255), song_number int, file varchar(255), song_order int, album_id int, FOREIGN KEY (album_id) REFERENCES albums(id) ON DELETE CASCADE)")
	s.songDB = sqlite.NewSongDB(db)
}

func (s *SongDBTestSuite) TearDownSuite() {
	defer s.db.Close()
	s.db.Exec("Drop table songs")
	s.db.Exec("Drop table albums")
}

func TestSongDBTestSuite(t *testing.T) {
	suite.Run(t, new(SongDBTestSuite))
}

func (s *SongDBTestSuite) TestCreateSong() {
	song, _ := domain.NewSong(dto.CreateSongDTO{
		FileName: "1 - test",
		AlbumID:  "album_id",
	})

	err := s.songDB.Create(song, nil)
	s.Nil(err)
	s.NotNil(song.ID)
	s.Equal("test", song.Title)
	s.Equal(1, song.Number)
	s.Equal("1 - test", song.File)
	s.Equal("album_id", song.AlbumID)

	song2, _ := domain.NewSong(dto.CreateSongDTO{
		FileName: "test",
		AlbumID:  "album_id",
	})

	err = s.songDB.Create(song2, nil)
	s.Nil(err)
	s.NotNil(song.ID)
	s.Equal("test", song2.Title)
	s.Equal(0, song2.Number)
	s.Equal("test", song2.File)
	s.Equal("album_id", song2.AlbumID)
}

func (s *SongDBTestSuite) TestFindAll() {
	album, _ := domain.NewAlbum(dto.CreateAlbumDTO{
		Title:    "title album",
		ImageUrl: "url",
	})

	err := s.albumDB.Create(album)
	s.Nil(err)

	song, _ := domain.NewSong(dto.CreateSongDTO{
		FileName: "1 - test",
		AlbumID:  album.ID,
	})

	song2, _ := domain.NewSong(dto.CreateSongDTO{
		FileName: "1 - test",
		AlbumID:  album.ID,
	})

	err = s.songDB.Create(song, nil)
	s.Nil(err)

	err = s.songDB.Create(song2, nil)
	s.Nil(err)

	songDTOS := []dto.SongWithAlbumDTO{
		{
			ID:            song.ID,
			Title:         song.Title,
			File:          song.File,
			SongNumber:    song.Number,
			Order:         1,
			AlbumID:       album.ID,
			AlbumName:     album.Title,
			AlbumImageURL: album.ImageUrl,
		},
		{
			ID:            song2.ID,
			Title:         song2.Title,
			File:          song2.File,
			SongNumber:    song2.Number,
			Order:         2,
			AlbumID:       album.ID,
			AlbumName:     album.Title,
			AlbumImageURL: album.ImageUrl,
		},
	}

	sg, err := s.songDB.FindAll()
	s.Nil(err)
	s.Equal(&songDTOS, sg)
}

func (s *SongDBTestSuite) TestUpdate() {
	song, _ := domain.NewSong(dto.CreateSongDTO{
		FileName: "1 - test",
		AlbumID:  "album test",
	})

	err := s.songDB.Create(song, nil)
	s.Nil(err)

	song.Title = "new title"

	err = s.songDB.Update(song)
	s.Nil(err)

	response, err := s.songDB.FindByID(song.ID)
	s.Nil(err)

	s.Equal("new title", response.Title)
}

func (s *SongDBTestSuite) TestUpdateOrder() {
	song, _ := domain.NewSong(dto.CreateSongDTO{
		FileName: "1 - test",
		AlbumID:  "album test",
	})

	song2, _ := domain.NewSong(dto.CreateSongDTO{
		FileName: "1 - test",
		AlbumID:  "album test",
	})

	song3, _ := domain.NewSong(dto.CreateSongDTO{
		FileName: "1 - test",
		AlbumID:  "album test",
	})

	err := s.songDB.Create(song, nil)
	s.Nil(err)

	err = s.songDB.Create(song2, nil)
	s.Nil(err)

	err = s.songDB.Create(song3, nil)
	s.Nil(err)

	dtos := []dto.SongUpdateOrderDTO{
		{
			ID:    song3.ID,
			Order: 1,
		},
		{
			ID:    song.ID,
			Order: 2,
		},
		{
			ID:    song2.ID,
			Order: 3,
		},
	}

	err = s.songDB.UpdateOrder(&dtos)
	s.Nil(err)

	response1, err := s.songDB.FindByID(song.ID)
	s.Nil(err)

	response2, err := s.songDB.FindByID(song2.ID)
	s.Nil(err)

	response3, err := s.songDB.FindByID(song3.ID)
	s.Nil(err)

	s.Equal(1, response3.Order)
	s.Equal(2, response1.Order)
	s.Equal(3, response2.Order)
}

func (s *SongDBTestSuite) TestDelete() {
	song, _ := domain.NewSong(dto.CreateSongDTO{
		FileName: "1 - test",
		AlbumID:  "album test",
	})

	err := s.songDB.Create(song, nil)
	s.Nil(err)

	err = s.songDB.Delete(song.ID)
	s.Nil(err)

	_, err = s.songDB.FindByID(song.ID)
	s.NotNil(err)
}
