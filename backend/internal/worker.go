package internal

import (
	"errors"
	"traceroute-optimised/internal/domain"
)

func Worker(workerID int, jobs chan domain.Job) {

	for job := range jobs {
		select {
		case <-job.Ctx.Done():
			job.Resp <- domain.Result{Status: 499, Error: errors.New("client cancelled")}
			continue
		default:
		}
		res := job.Exec()

		job.Resp <- res
	}
}
