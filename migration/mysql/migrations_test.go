package mysql

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/suite"
)

type MigrationsTests struct {
	suite.Suite
}

func TestMigrations(t *testing.T) {
	suite.Run(t, &MigrationsTests{})
}

func (s *MigrationsTests) TestSortable() {
	migrations := Migrations{
		New("2", "up", "down"),
		New("1", "up", "down"),
		New("b", "up", "down"),
		New("a", "up", "down"),
		New("B", "up", "down"),
		New("A", "up", "down"),
	}
	sort.Sort(migrations)
	s.Equal(migrations[0].Name, "1")
	s.Equal(migrations[1].Name, "2")
	s.Equal(migrations[2].Name, "A")
	s.Equal(migrations[3].Name, "B")
	s.Equal(migrations[4].Name, "a")
	s.Equal(migrations[5].Name, "b")
}
