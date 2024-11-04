package main

import (
	"fmt"
)

func main() {
	wp := NewWorkerPool()
	wp.StartProcessingResults()

	wp.AddWorker()

	for i := 0; i < 10; i++ {
		wp.AddJob(fmt.Sprintf("Job %d", i+1))
	}
	wp.DeleteWorker()
	wp.AddWorkers(2)

	for i := 0; i < 10; i++ {
		wp.AddJob(fmt.Sprintf("Job %d", i+11))
	}

	wp.DeleteWorkers(3)

	for i := 0; i < 10; i++ {
		wp.AddJob(fmt.Sprintf("Job %d", i+21))
	}

	wp.Stop()

}
