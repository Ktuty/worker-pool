package main

import (
	"fmt"
	"sync"
)

type Worker struct {
	id   int
	job  chan string
	stop chan struct{}
	wg   sync.WaitGroup
}

func NewWorker(id int) *Worker {
	return &Worker{
		id:   id,
		job:  make(chan string),
		stop: make(chan struct{}),
	}
}

type WorkerPool struct {
	Workers []*Worker
	res     chan string
	jobs    chan string
	mu      sync.Mutex
}

func NewWorkerPool() *WorkerPool {
	return &WorkerPool{
		res:  make(chan string),
		jobs: make(chan string),
	}
}

func (wp *WorkerPool) AddWorker() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	worker := NewWorker(len(wp.Workers))
	wp.Workers = append(wp.Workers, worker)

	worker.wg.Add(1)
	go wp.start(worker)
}

func (wp *WorkerPool) AddWorkers(count int) {
	for i := 0; i < count; i++ {
		wp.AddWorker()
	}
}

func (wp *WorkerPool) DeleteWorkers(count int) {
	for i := 0; i < count; i++ {
		wp.DeleteWorker()
	}
}

func (wp *WorkerPool) DeleteWorker() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if len(wp.Workers)-1 > 0 {
		worker := wp.Workers[len(wp.Workers)-1]
		close(worker.stop)
		wp.Workers = wp.Workers[:len(wp.Workers)-1]
		fmt.Printf("Delete worker №%d\n", len(wp.Workers)+1)
		worker.wg.Wait()
	} else {
		fmt.Printf("you can't delete the only one worker (№%d)\n", len(wp.Workers))
	}
}

func (wp *WorkerPool) start(worker *Worker) {
	defer worker.wg.Done()

	fmt.Printf("Start processing №%d\n", worker.id+1)
	for {
		select {
		case job, ok := <-wp.jobs:
			if !ok {
				return
			}
			wp.res <- fmt.Sprintf("Worker №%d: %s", worker.id+1, job)
		case <-worker.stop:
			return
		}
	}
}

func (wp *WorkerPool) AddJob(job string) {
	wp.jobs <- job
}

func (wp *WorkerPool) Stop() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	close(wp.jobs)
	for _, worker := range wp.Workers {
		worker.wg.Wait()
	}
	close(wp.res)
}

func (wp *WorkerPool) StartProcessingResults() {
	go func() {
		for result := range wp.res {
			fmt.Println(result)
		}
	}()
}
