package solution

import (
	"errors"
	"time"
)

// WorkerPool errors, do not change!
var (
	ErrBadParams  = errors.New("bad params")
	ErrBadTask    = errors.New("bad task")
	ErrNotStarted = errors.New("not started")
)

// WorkerPool represents a pool of goroutines.
type WorkerPool struct {
	free    int
	running bool
	jobs    chan Task
	results chan error
}

// Task to be computed by the WorkerPool.
type Task func() error

// NewWorkerPool creates a new pool with a given size.
func NewWorkerPool(size int) (*WorkerPool, error) {
	if size > 10 {
		return nil, ErrBadParams
	}

	workerPool := WorkerPool{
		free:    size,
		running: false,
		jobs:    make(chan Task),
		results: make(chan error),
	}
	return &workerPool, nil
}

// Results returns channel of non-nil errors.
func (wp *WorkerPool) Results() <-chan error {
	return wp.results
}

// Run will start jobs(goroutines) for tasks computation.
func (wp *WorkerPool) Run() {
	if wp.running {
		return
	}
	wp.running = true
	go func() {
		for true {
			if wp.free > 0 {
				job := <-wp.jobs
				wp.free = wp.free - 1
				go func() {
					err := job()
					wp.free = wp.free + 1
					if err != nil {
						wp.results <- err
					}
				}()
			}
			time.Sleep(1 * time.Millisecond)
		}
	}()
}

// AddTask will add a task to the worker pool queue.
func (wp *WorkerPool) AddTask(task Task) error {
	if task == nil {
		return ErrBadTask
	}
	if !wp.running {
		return ErrNotStarted
	}
	wp.jobs <- task
	return nil
}
