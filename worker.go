package main

import (
	"sync"
)

type Worker interface {
	Work() error
}

type Runner struct {
	StopChan chan struct{}
}

func NewRunner() *Runner {
	r := &Runner{}
	r.StopChan = make(chan struct{}, 1)
	return r
}

func (r *Runner) Run(workers []Worker) {
	var wg sync.WaitGroup
	stopNotifyChan := make(chan struct{}, len(workers)) // for non blocking notify

	for _, worker := range workers {
		wg.Add(1)
		go func(worker Worker) {
			defer wg.Done()
			if err := worker.Work(); err != nil {
				Info("worker err ", worker, " ", err)
				stopNotifyChan <- struct{}{}
			}
		}(worker)
	}

	<-stopNotifyChan
	Info("close stopchan")
	close(r.StopChan)

	wg.Wait()
	Info("all workers are died")
}
