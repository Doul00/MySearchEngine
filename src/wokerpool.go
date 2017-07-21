package main

import (
	"bufio"
)

type Worker struct {
	queryChan  chan Query
	workerChan chan *Worker
	finishChan chan bool
}

type Query struct {
	index *Index
	query string
}

func (w *Worker) run() {
	go func() {
		for {
			w.workerChan <- w

			select {
			case query := <-w.queryChan:
				results := search(query.query, query.index)
				printResults(results)

			case <-w.finishChan:
				return
			}
		}
	}()
}

func (w *Worker) stop() {
	go func() {
		w.finishChan <- true
	}()
}

func createWorker(workerChan chan *Worker) *Worker {
	return &Worker{
		queryChan:  make(chan Query),
		workerChan: workerChan,
		finishChan: make(chan bool, 1)}
}

func startDispatching(scanner *bufio.Scanner, indexes []*Index, nbGo int) {
	workerQueue := make(chan *Worker)

	for i := 0; i < nbGo; i++ {
		worker := createWorker(workerQueue)
		worker.run()
	}

	for scanner.Scan() {
		query := scanner.Text()

		for _, index := range indexes {
			worker := <-workerQueue
			worker.queryChan <- Query{index, query}
		}
	}

	for i := 0; i < nbGo; i++ {
		worker := <-workerQueue
		worker.stop()
	}
}
