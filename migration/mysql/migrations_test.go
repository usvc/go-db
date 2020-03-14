package mysql

import (
	"fmt"
	"sort"
	"strconv"
	"testing"

	"github.com/stretchr/testify/suite"
)

type MigrationsTests struct {
	suite.Suite
}

func TestMigrations(t *testing.T) {
	suite.Run(t, &MigrationsTests{})
}

func (s *MigrationsTests) TestSortable_names() {
	up := "up"
	down := "down"
	migrations := Migrations{}
	testOrder := []string{
		"20200302030405_b_000005",
		"20200302030405_a_000004",
		"20200302030405_1_000000",
		"20200302030405_B_000003",
		"20200302030405_A_000002",
		"20200302030405_2_000001",
	}
	for i := 0; i < len(testOrder); i++ {
		migrations = append(migrations, New(testOrder[i], up, down))
	}
	sort.Sort(migrations)
	for i := 0; i < len(migrations); i++ {
		s.Contains(migrations[i].Name, fmt.Sprintf("00000%v", strconv.Itoa(i)))
	}
}

func (s *MigrationsTests) TestSortable_dates() {
	up := "up"
	down := "down"
	migrations := Migrations{}
	testOrder := []string{
		"20200202030405_migration_000005",
		"20200102040405_migration_000003",
		"20200102030406_migration_000001",
		"20200103030405_migration_000004",
		"20200102030505_migration_000002",
		"20200102030405_migration_000000",
	}
	for i := 0; i < len(testOrder); i++ {
		migrations = append(migrations, New(testOrder[i], up, down))
	}
	sort.Sort(migrations)
	for i := 0; i < len(migrations); i++ {
		s.Contains(migrations[i].Name, fmt.Sprintf("00000%v", strconv.Itoa(i)))
	}
}
