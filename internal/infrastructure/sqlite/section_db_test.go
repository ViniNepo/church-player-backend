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

type SectionDBTestSuite struct {
	suite.Suite
	db        *sql.DB
	sectionDB *sqlite.SectionDB
	momentDB  *sqlite.MomentDB
}

func (s *SectionDBTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.db = db

	db.Exec("CREATE table sections (id varchar(255), section_name varchar(255), section_order int, worship_id varchar(255))")
	s.sectionDB = sqlite.NewSectionDB(db)
	db.Exec("CREATE table moments (id varchar(255), title varchar(255), song_order int, song_id int, section_id int,  FOREIGN KEY (section_id) REFERENCES sections(id) ON DELETE CASCADE)")
	s.momentDB = sqlite.NewMomentDB(db)
}

func (s *SectionDBTestSuite) TearDownSuite() {
	defer s.db.Close()
	s.db.Exec("Drop table sections")
	s.db.Exec("Drop table moments")
}

func TestSectionDBTestSuite(t *testing.T) {
	suite.Run(t, new(SectionDBTestSuite))
}

func (s *SectionDBTestSuite) TestCreate() {
	section, _ := domain.NewSection(dto.CreateSectionDTO{
		SectionName: "name",
		WorshipID:   "id",
	})

	err := s.sectionDB.Create(section)
	s.Nil(err)

	section2, _ := domain.NewSection(dto.CreateSectionDTO{
		SectionName: "name",
		WorshipID:   "id",
	})

	err = s.sectionDB.Create(section2)
	s.Nil(err)
}

func (s *SectionDBTestSuite) TestFindAllByWorshipID() {
	section, _ := domain.NewSection(dto.CreateSectionDTO{
		SectionName: "name",
		WorshipID:   "worship_id",
	})

	err := s.sectionDB.Create(section)
	s.Nil(err)

	section2, _ := domain.NewSection(dto.CreateSectionDTO{
		SectionName: "name",
		WorshipID:   "worship_id",
	})

	err = s.sectionDB.Create(section2)
	s.Nil(err)

	section3, _ := domain.NewSection(dto.CreateSectionDTO{
		SectionName: "name",
		WorshipID:   "worship_id2",
	})

	err = s.sectionDB.Create(section3)
	s.Nil(err)

	response, err := s.sectionDB.FindAllByWorshipID("worship_id")
	s.Nil(err)
	s.Equal(2, len(*response))
}

func (s *SectionDBTestSuite) TestUpdate() {
	section, _ := domain.NewSection(dto.CreateSectionDTO{
		SectionName: "name",
		WorshipID:   "worship_id2",
	})

	err := s.sectionDB.Create(section)
	s.Nil(err)

	section.SectionName = "new Title"

	err = s.sectionDB.Update(section)
	s.Nil(err)
}

func (s *SectionDBTestSuite) TestDeleteSectionByWorshipId() {
	section, _ := domain.NewSection(dto.CreateSectionDTO{
		SectionName: "name",
		WorshipID:   "worship_id",
	})

	err := s.sectionDB.Create(section)
	s.Nil(err)

	section2, _ := domain.NewSection(dto.CreateSectionDTO{
		SectionName: "name",
		WorshipID:   "worship_id",
	})

	err = s.sectionDB.Create(section2)
	s.Nil(err)

	err = s.sectionDB.Delete(section.ID)
	s.Nil(err)

	response, err := s.sectionDB.FindAllByWorshipID("worship_id")
	s.Nil(err)
	s.Equal(1, len(*response))
}
