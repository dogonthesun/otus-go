package hw05parallelexecution

import (
	"errors"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	taskCh := make(chan Task)
	errCh := make(chan error)

	if m <= 0 {
		m = len(tasks) + 1
	}

	for range n {
		go func() {
			for task := range taskCh {
				errCh <- task()
			}
		}()
	}

	enqueued, finished, errors := 0, 0, 0

	for enqueued < len(tasks) && errors < m {
		select {
		case err := <-errCh:
			finished++
			if err != nil {
				errors++
			}
		case taskCh <- tasks[enqueued]:
			enqueued++
		}
	}

	close(taskCh)

	for range enqueued - finished {
		if err := <-errCh; err != nil {
			errors++
		}
	}

	if errors >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
