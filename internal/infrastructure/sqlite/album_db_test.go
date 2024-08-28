package sqlite_test

import (
	"database/sql"
	"testing"

	"github.com/ViniNepo/church-player-backend/internal/domain"
	"github.com/ViniNepo/church-player-backend/internal/domain/dto"
	"github.com/ViniNepo/church-player-backend/internal/infrastructure/sqlite"
	"github.com/stretchr/testify/suite"

	_ "github.com/mattn/go-sqlite3"
)

type AlbumDBTestSuite struct {
	suite.Suite
	db      *sql.DB
	albumDB *sqlite.AlbumDB
	songDB  *sqlite.SongDB
}

func (s *AlbumDBTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.db = db

	db.Exec("CREATE table albums (id varchar(255), title varchar(255), image_url varchar(255), album_order int)")
	s.albumDB = sqlite.NewAlbumDB(db)
	db.Exec("CREATE table songs (id varchar(255), title varchar(255), song_number int, file varchar(255), song_order int, album_id int, FOREIGN KEY (album_id) REFERENCES albums(id) ON DELETE CASCADE)")
	s.songDB = sqlite.NewSongDB(db)
}

func (s *AlbumDBTestSuite) TearDownSuite() {
	defer s.db.Close()
	s.db.Exec("Drop table albums")
	s.db.Exec("Drop table songs")
}

func TestAlbumDBTestSuite(t *testing.T) {
	suite.Run(t, new(AlbumDBTestSuite))
}

func (s *AlbumDBTestSuite) TestCreateAlbum() {
	album, _ := domain.NewAlbum(dto.CreateAlbumDTO{
		Title:    "Title",
		ImageUrl: "",
	})

	err := s.albumDB.Create(album)
	s.Nil(err)

	album2, _ := domain.NewAlbum(dto.CreateAlbumDTO{
		Title:    "Title",
		ImageUrl: "Url",
	})

	err = s.albumDB.Create(album2)
	s.Nil(err)
}

func (s *AlbumDBTestSuite) TestCreateAlbumWithoutTitle() {
	album, err := domain.NewAlbum(dto.CreateAlbumDTO{
		Title:    "",
		ImageUrl: "",
	})
	s.Nil(album)
	s.NotNil(err)

	err = s.albumDB.Create(album)
	s.NotNil(err)
}

func (s *AlbumDBTestSuite) TestAlbumById() {
	album, _ := domain.NewAlbum(dto.CreateAlbumDTO{
		Title:    "Title",
		ImageUrl: "",
	})

	err := s.albumDB.Create(album)
	s.Nil(err)

	response, err := s.albumDB.FindByID(album.ID)
	s.Nil(err)
	s.NotNil(response)
	s.Equal("Title", response.Title)
	s.Equal("", response.ImageUrl)

	_, err = s.albumDB.FindByID("")
	s.NotNil(err)
}

func (s *AlbumDBTestSuite) TestFindAll() {
	album, _ := domain.NewAlbum(dto.CreateAlbumDTO{
		Title:    "Title",
		ImageUrl: "",
	})

	err := s.albumDB.Create(album)
	s.Nil(err)

	album2, _ := domain.NewAlbum(dto.CreateAlbumDTO{
		Title:    "Title 2",
		ImageUrl: "",
	})

	err = s.albumDB.Create(album2)
	s.Nil(err)

	response, err := s.albumDB.FindAll()
	s.Nil(err)
	s.NotNil(response)
}

func (s *AlbumDBTestSuite) TestFindAllWithoutResult() {
	_, err := s.albumDB.FindAll()
	s.Nil(err)
}

func (s *AlbumDBTestSuite) TestUpdateAlbum() {
	album, _ := domain.NewAlbum(dto.CreateAlbumDTO{
		Title:    "Title",
		ImageUrl: "",
	})

	err := s.albumDB.Create(album)
	s.Nil(err)

	album.Title = "new Title"

	err = s.albumDB.Update(album)
	s.Nil(err)

	response, err := s.albumDB.FindByID(album.ID)
	s.Nil(err)

	s.NotNil(response)
	s.Equal("new Title", response.Title)
}

func (s *AlbumDBTestSuite) TestDeleteAlbumById() {
	album, _ := domain.NewAlbum(dto.CreateAlbumDTO{
		Title:    "Title",
		ImageUrl: "",
	})

	song, _ := domain.NewSong(dto.CreateSongDTO{
		FileName: "1 - test",
		AlbumID:  album.ID,
	})

	err := s.albumDB.Create(album)
	s.Nil(err)

	err = s.songDB.Create(song)
	s.Nil(err)
	s.NotNil(song.ID)

	err = s.albumDB.Delete(album.ID)
	s.Nil(err)

	_, err = s.albumDB.FindByID(album.ID)
	s.NotNil(err)

	_, err = s.songDB.FindByID(song.ID)
	s.NotNil(err)
}
