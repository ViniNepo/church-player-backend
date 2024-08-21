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

	db.Exec("CREATE table moments (id varchar(255), title varchar(255), section_name varchar(255), section_order int, song_order int, song_id int, worship_id int,  FOREIGN KEY (worship_id) REFERENCES worships(id) ON DELETE CASCADE)")
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
		Title:        "new moment",
		SectionName:  "inicio",
		SectionOrder: 1,
		SongOrder:    1,
		WorshipID:    "id",
	})

	err := s.momentDB.Create(moment)
	s.Nil(err)
}

func (s *MomentDBTestSuite) TestFindByID() {
	moment, _ := domain.NewMoment(dto.CreateMomentDTO{
		Title:        "new moment",
		SectionName:  "inicio",
		SectionOrder: 1,
		SongOrder:    1,
		WorshipID:    "id",
	})

	err := s.momentDB.Create(moment)
	s.Nil(err)

	m, err := s.momentDB.FindByID(moment.ID)
	s.Nil(err)
	s.Equal(moment, m)
}

func (s *MomentDBTestSuite) TestFindAllByWorshipID() {
	moment1, _ := domain.NewMoment(dto.CreateMomentDTO{
		Title:        "new moment",
		SectionName:  "inicio",
		SectionOrder: 1,
		SongOrder:    1,
		WorshipID:    "id",
	})

	moment2, _ := domain.NewMoment(dto.CreateMomentDTO{
		Title:        "new moment",
		SectionName:  "inicio",
		SectionOrder: 1,
		SongOrder:    1,
		WorshipID:    "id",
	})

	err := s.momentDB.Create(moment1)
	s.Nil(err)

	err = s.momentDB.Create(moment2)
	s.Nil(err)

	moments, err := s.momentDB.FindAllByWorshipID("id")
	s.Nil(err)
	s.NotNil(moments)
}

func (s *MomentDBTestSuite) TestUpdate() {
	moment, _ := domain.NewMoment(dto.CreateMomentDTO{
		Title:        "new moment",
		SectionName:  "inicio",
		SectionOrder: 1,
		SongOrder:    1,
		WorshipID:    "id",
	})

	err := s.momentDB.Create(moment)
	s.Nil(err)

	moment.Title = "new title"

	err = s.momentDB.Update(moment)
	s.Nil(err)

	m, err := s.momentDB.FindByID(moment.ID)
	s.Nil(err)
	s.Equal("new title", m.Title)
}

func (s *MomentDBTestSuite) TestUpdateOrder() {
	moment, _ := domain.NewMoment(dto.CreateMomentDTO{
		Title:        "new moment",
		SectionName:  "inicio",
		SectionOrder: 1,
		SongOrder:    1,
		WorshipID:    "test",
	})

	moment2, _ := domain.NewMoment(dto.CreateMomentDTO{
		Title:        "new moment",
		SectionName:  "inicio",
		SectionOrder: 1,
		SongOrder:    2,
		WorshipID:    "test",
	})

	moment3, _ := domain.NewMoment(dto.CreateMomentDTO{
		Title:        "new moment",
		SectionName:  "inicio",
		SectionOrder: 1,
		SongOrder:    3,
		WorshipID:    "test",
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

	m, err := s.momentDB.FindAllByWorshipID("test")
	s.Nil(err)
	s.NotNil(m)
}

func (s *MomentDBTestSuite) TestDelete() {
	moment, _ := domain.NewMoment(dto.CreateMomentDTO{
		Title:        "new moment",
		SectionName:  "inicio",
		SectionOrder: 1,
		SongOrder:    1,
		WorshipID:    "id",
	})

	err := s.momentDB.Create(moment)
	s.Nil(err)

	err = s.momentDB.Delete(moment.ID)
	s.Nil(err)
}
