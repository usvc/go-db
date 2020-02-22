package db

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/DATA-DOG/go-sqlmock"
)

type DBTests struct {
	suite.Suite
}

func TestDB(t *testing.T) {
	suite.Run(t, &DBTests{})
}

func (s *DBTests) TestCheck_ping() {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	s.Nil(err)
	defer db.Close()
	Import(db, "__check_ping")
	fmt.Println("hi")
	mock.ExpectPing()
	Check("__check_ping")
	s.Nil(mock.ExpectationsWereMet())
}

func (s *DBTests) TestCheck_pingReturnsError() {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	s.Nil(err)
	defer db.Close()
	Import(db, "__check_ping_returns_error")
	fmt.Println("hi")
	mock.ExpectPing().WillReturnError(fmt.Errorf("this is expected"))
	err = Check("__check_ping_returns_error")
	s.Equal("this is expected", err.Error())
	s.Nil(mock.ExpectationsWereMet())
}

func (s *DBTests) TestCheck_notFound() {
	err := Check("__check_not_found")
	s.Contains(err.Error(), "does not exist")
}


func (s *DBTests) TestCloseAll() {
	db, mock, err := sqlmock.New()
	s.Nil(err)
	defer db.Close()
	mock.ExpectClose()

	db1, mock1, err := sqlmock.New()
	s.Nil(err)
	defer db1.Close()
	mock1.ExpectClose()
	
	s.Nil(Import(db, "__close_all_0"))
	s.Nil(Import(db1, "__close_all_1"))
	errs := CloseAll()
	s.Nil(errs)
}

func (s *DBTests) TestGet() {
	db, _, err := sqlmock.New()
	s.Nil(err)
	defer db.Close()
	err = Import(db, "__test_get")
	s.NotNil(Get("__test_get"))
	
}

func (s *DBTests) TestImport() {
	db, _, err := sqlmock.New()
	s.Nil(err)
	defer db.Close()
	err = Import(db, "__test_import")
	s.Nil(err)
}

func (s *DBTests) TestInit_closesExisting() {
	db, mock, err := sqlmock.New()
	s.Nil(err)
	mock.ExpectClose()
	Import(db, "__test_init_closes_existing")
	err = Init(Options{
		ConnectionName: "__test_init_closes_existing",
	})
	s.Nil(err)
	s.Nil(mock.ExpectationsWereMet())
}