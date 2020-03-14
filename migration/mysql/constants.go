package mysql

import "errors"

const (
	StatusRollingBack = "rolling back"
	StatusApplying    = "applying"
	StatusApplied     = "applied"
)

var (
	NoErrAlreadyApplied = errors.New("migration has already been applied")
	NoErrDoesNotExist   = errors.New("migration does not exist")
)
