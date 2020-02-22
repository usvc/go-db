package db

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type OptionsTests struct {
	suite.Suite
}

func TestOptions(t *testing.T) {
	suite.Run(t, &OptionsTests{})
}

func (s *OptionsTests) TestAssignDefaults() {
	options := Options{}
	options.AssignDefaults()
	s.Equal(DefaultConnectionName, options.ConnectionName)
	s.Equal(DefaultDriver, options.Driver)
}
