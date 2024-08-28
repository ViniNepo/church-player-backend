package sqlite_test

import (
	"database/sql"
	"testing"

	"github.com/ViniNepo/church-player-backend/internal/domain"
	"github.com/ViniNepo/church-player-backend/internal/domain/dto"
	"github.com/ViniNepo/church-player-backend/internal/infrastructure/sqlite"
	"github.com/stretchr/testify/suite"
)

type MomentDBTestSuite struct {
	suite.Suite
	db       *sql.DB
	momentDB *sqlite.MomentDB
}

func (s *MomentDBTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.db = db

	db.Exec("CREATE table moments (id varchar(255), title varchar(255), song_order int, song_id int, section_id int,  FOREIGN KEY (section_id) REFERENCES sections(id) ON DELETE CASCADE)")
	s.momentDB = sqlite.NewMomentDB(db)
}

func (s *MomentDBTestSuite) TearDownSuite() {
	defer s.db.Close()
	s.db.Exec("Drop table moments")
}

func TestMomentDBTestSuite(t *testing.T) {
	suite.Run(t, new(MomentDBTestSuite))
}

func (s *MomentDBTestSuite) TestCreate() {
	moment, _ := domain.NewMoment(dto.CreateMomentDTO{
		Title:     "new moment",
		SectionID: "id",
	})

	err := s.momentDB.Create(moment)
	s.Nil(err)
}

func (s *MomentDBTestSuite) TestFindByID() {
	moment, _ := domain.NewMoment(dto.CreateMomentDTO{
		Title:     "new moment",
		SectionID: "section id",
	})

	err := s.momentDB.Create(moment)
	s.Nil(err)

	m, err := s.momentDB.FindByID(moment.ID)
	s.Nil(err)
	s.Equal(1, m.SongOrder)
}

func (s *MomentDBTestSuite) TestUpdate() {
	moment, _ := domain.NewMoment(dto.CreateMomentDTO{
		Title:     "new moment",
		SectionID: "id",
	})

	err := s.momentDB.Create(moment)
	s.Nil(err)

	moment.Title = "new title"
	moment.SongID = "song id"

	err = s.momentDB.Update(moment)
	s.Nil(err)

	m, err := s.momentDB.FindByID(moment.ID)
	s.Nil(err)
	s.Equal("new title", m.Title)
	s.Equal("song id", m.SongID)
}

func (s *MomentDBTestSuite) TestUpdateOrder() {
	moment, _ := domain.NewMoment(dto.CreateMomentDTO{
		Title:     "new moment",
		SectionID: "test",
	})

	moment2, _ := domain.NewMoment(dto.CreateMomentDTO{
		Title:     "new moment 2",
		SectionID: "test",
	})

	moment3, _ := domain.NewMoment(dto.CreateMomentDTO{
		Title:     "new moment 3",
		SectionID: "test",
	})

	err := s.momentDB.Create(moment)
	s.Nil(err)

	err = s.momentDB.Create(moment2)
	s.Nil(err)

	err = s.momentDB.Create(moment3)
	s.Nil(err)

	dtos := []dto.MomentOrderDTO{
		{
			ID:        moment3.ID,
			SongOrder: 1,
		},
		{
			ID:        moment.ID,
			SongOrder: 2,
		},
		{
			ID:        moment2.ID,
			SongOrder: 3,
		},
	}

	err = s.momentDB.UpdateOrder(&dtos)
	s.Nil(err)

	m, err := s.momentDB.FindByID(moment.ID)
	s.Nil(err)
	s.Equal("new moment", m.Title)
	s.Equal(2, m.SongOrder)

	m, err = s.momentDB.FindByID(moment2.ID)
	s.Nil(err)
	s.Equal("new moment 2", m.Title)
	s.Equal(3, m.SongOrder)

	m, err = s.momentDB.FindByID(moment3.ID)
	s.Nil(err)
	s.Equal("new moment 3", m.Title)
	s.Equal(1, m.SongOrder)
}

func (s *MomentDBTestSuite) TestDelete() {
	moment, _ := domain.NewMoment(dto.CreateMomentDTO{
		Title:     "new moment",
		SectionID: "test",
	})

	err := s.momentDB.Create(moment)
	s.Nil(err)

	err = s.momentDB.Delete(moment.ID)
	s.Nil(err)
}
