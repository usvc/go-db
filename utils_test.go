package db

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type UtilsTests struct {
	suite.Suite
}

func TestUtils(t *testing.T) {
	suite.Run(t, &UtilsTests{})
}

func (s *UtilsTests) Test_generateDSN() {
	dsn := generateDSN(Options{
		Username: "_username",
		Password: "_password",
		Database: "_database",
		Hostname: "_host._name",
		Port:     98765,
		Params: map[string]string{
			"bool":   "true",
			"int":    "1234",
			"float":  "3.14",
			"string": "string",
		},
	})
	s.Equal("_username:_password@tcp(_host._name:98765)/_database?bool=true&float=3.14&int=1234&parseTime=true&string=string", dsn)
}
