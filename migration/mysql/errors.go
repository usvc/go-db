package mysql

import "errors"

const (
	StatusRollingBack = "rolling back"
	StatusApplying    = "applying"
	StatusApplied     = "applied"
)

var (
	NoErrAlreadyApplied = errors.New("migration has already been applied")
)
