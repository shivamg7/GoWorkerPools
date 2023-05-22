package main

import (
	"AspireLoans/solution"
	"errors"
	"fmt"
	"time"
)

func someLongRunningTask() error {
	fmt.Println("long running task started")
	time.Sleep(20 * time.Second)
	fmt.Println("completed long running task")
	return nil
}

func someShortTask() error {
	fmt.Println("short running task started")
	time.Sleep(2 * time.Second)
	fmt.Println("completed short running task")
	return nil
}

func errorTask() error {
	return errors.New("something went wrong")
}

func main() {
	// create a pool with three workers,
	// which means three tasks can be done in parallel.
	wp, err := solution.NewWorkerPool(3)
	if err != nil {
		panic("should be ok")
	}

	wp.Run() // run the pool
	wp.Run()

	for i := 0; i < 4; i++ {
		err := wp.AddTask(someLongRunningTask)
		if err != nil {
			panic("should fail only on nil task")
		}
	}

	// add four tasks that will be processed in a short time
	for i := 0; i < 4; i++ {
		err := wp.AddTask(someShortTask)
		if err != nil {
			panic("should fail only on nil task")
		}
	}

	err = wp.AddTask(errorTask)
	if err != nil {
		panic("should fail only on nil task")
	}
	// after a short time all of the short tasks should be processed here
	// even if we have only three workers and five tasks were submitted
	// just because `someLongRunningTask` is still in process
	// but the four short tasks are already finished.

	// get the error results from worker pool
	for err := range wp.Results() {
		fmt.Printf("got err: %v\n", err)
	}

	<-make(chan int)
}
