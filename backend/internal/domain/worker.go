package domain

import (
	"context"
)

type JobExecFunc func() Result
type Job struct {
	ID   int64
	Ctx  context.Context
	Resp chan Result
	Exec JobExecFunc
}

type Result struct {
	Status int
	Body   any
	Error  error
}
