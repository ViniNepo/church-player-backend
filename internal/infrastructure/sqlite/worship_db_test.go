package sqlite_test

import (
	"database/sql"
	"testing"

	"github.com/ViniNepo/church-player-backend/internal/domain"
	"github.com/ViniNepo/church-player-backend/internal/domain/dto"
	"github.com/ViniNepo/church-player-backend/internal/infrastructure/sqlite"
	"github.com/stretchr/testify/suite"
)

type WorshipDBTestSuite struct {
	suite.Suite
	db        *sql.DB
	worshipDB *sqlite.WorshipDB
	momentDB  *sqlite.MomentDB
	albumDB   *sqlite.AlbumDB
	songDB    *sqlite.SongDB
}

func (s *WorshipDBTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.db = db

	db.Exec("CREATE table worships (id varchar(255), title varchar(255), image_url varchar(255), worship_order int)")
	s.worshipDB = sqlite.NewWorshipDB(db)
}

func (s *WorshipDBTestSuite) TearDownSuite() {
	defer s.db.Close()
	s.db.Exec("Drop table worships")
}

func TestWorshipDBTestSuite(t *testing.T) {
	suite.Run(t, new(WorshipDBTestSuite))
}

func (s *WorshipDBTestSuite) TestCreateWorship() {
	worship, _ := domain.NewWorship(dto.CreateWorshipDTO{
		Title:    "Title",
		ImageUrl: "",
	})

	err := s.worshipDB.Create(worship)
	s.Nil(err)

	worship2, _ := domain.NewWorship(dto.CreateWorshipDTO{
		Title:    "Title",
		ImageUrl: "Url",
	})

	err = s.worshipDB.Create(worship2)
	s.Nil(err)
}

func (s *WorshipDBTestSuite) TestCreateWorshipWithoutTitle() {
	worship, err := domain.NewWorship(dto.CreateWorshipDTO{
		Title:    "",
		ImageUrl: "",
	})
	s.Nil(worship)
	s.NotNil(err)

	err = s.worshipDB.Create(worship)
	s.NotNil(err)
}

func (s *WorshipDBTestSuite) TestFindById() {
	worship, _ := domain.NewWorship(dto.CreateWorshipDTO{
		Title:    "Title",
		ImageUrl: "",
	})

	err := s.worshipDB.Create(worship)
	s.Nil(err)

	response, err := s.worshipDB.FindByID(worship.ID)
	s.Nil(err)
	s.NotNil(response)
	s.Equal("Title", response.Title)
	s.Equal("", response.ImageUrl)
}

func (s *WorshipDBTestSuite) TestFindAll() {
	worship, _ := domain.NewWorship(dto.CreateWorshipDTO{
		Title:    "Title",
		ImageUrl: "",
	})

	err := s.worshipDB.Create(worship)
	s.Nil(err)

	worship2, _ := domain.NewWorship(dto.CreateWorshipDTO{
		Title:    "Title 2",
		ImageUrl: "",
	})

	err = s.worshipDB.Create(worship2)
	s.Nil(err)

	response, err := s.worshipDB.FindAll()
	s.Nil(err)
	s.NotNil(response)
}

func (s *WorshipDBTestSuite) TestFindAllWithoutResult() {
	_, err := s.worshipDB.FindAll()
	s.Nil(err)
}

func (s *WorshipDBTestSuite) TestUpdateWorship() {
	worship, _ := domain.NewWorship(dto.CreateWorshipDTO{
		Title:    "Title",
		ImageUrl: "",
	})

	err := s.worshipDB.Create(worship)
	s.Nil(err)

	worship.Title = "new Title"

	err = s.worshipDB.Update(worship)
	s.Nil(err)

	response, err := s.worshipDB.FindByID(worship.ID)
	s.Nil(err)

	s.NotNil(response)
	s.Equal("new Title", response.Title)
}

func (s *WorshipDBTestSuite) TestDeleteWorshipById() {
	worship, _ := domain.NewWorship(dto.CreateWorshipDTO{
		Title:    "Title",
		ImageUrl: "",
	})

	err := s.worshipDB.Create(worship)
	s.Nil(err)

	err = s.worshipDB.Delete(worship.ID)
	s.Nil(err)

	_, err = s.worshipDB.FindByID(worship.ID)
	s.NotNil(err)
}
