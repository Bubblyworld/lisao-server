package tasks

import (
	"errors"
	"log"
	"sync"
	"time"
)

const workChannelCapacity = 20

type ProcessFunc func(interface{}) error

// A pool manages state of a pool of workers that read in tasks from a channel
// and process them in some way.
type pool struct {
	processFunc ProcessFunc
	poolLabel   string
	numWorkers  int

	stopChan  chan bool
	workChan  <-chan interface{}
	waitGroup *sync.WaitGroup

	isRunning bool
}

func NewPool(poolLabel string, numWorkers int, processFunc ProcessFunc,
	workChan <-chan interface{}) *pool {

	return &pool{
		processFunc: processFunc,
		poolLabel:   poolLabel,
		numWorkers:  numWorkers,

		stopChan:  make(chan bool),
		workChan:  workChan,
		waitGroup: &sync.WaitGroup{},

		isRunning: false,
	}
}

func (p *pool) Run() {
	p.isRunning = true

	var workerStopChans []chan bool
	for i := 0; i < p.numWorkers; i++ {
		workerStopChan := make(chan bool)
		workerStopChans = append(workerStopChans, workerStopChan)

		p.waitGroup.Add(1)
		go p.workForever(i, workerStopChan, p.waitGroup)
	}

	// Wait for stop signal - if it fires, stop the workers and return.
	for {
		select {
		case shouldStop := <-p.stopChan:
			if shouldStop {
				for _, workerStopChan := range workerStopChans {
					workerStopChan <- true
				}

				break
			}
		}

		time.Sleep(time.Second)
	}
}

func (p *pool) Stop() error {
	if !p.isRunning {
		return errors.New("pool: can't stop a pool that hasn't been started")
	}
	p.stopChan <- true
	p.waitGroup.Wait()
	p.isRunning = false
	return nil
}

func (p *pool) workForever(workerIndex int, stopChan <-chan bool,
	waitGroup *sync.WaitGroup) {

	log.Printf("Pool %s: starting worker %d.", p.poolLabel, workerIndex)
	defer func() {
		waitGroup.Done()
	}()

workForeverLoop:
	for {
		select {
		case shouldStop := <-stopChan:
			if shouldStop {
				log.Printf("Pool %s: stopping worker %d.", p.poolLabel, workerIndex)
				break workForeverLoop
			}

		case work := <-p.workChan:
			log.Printf("Pool %s: received work for worker %d.", p.poolLabel, workerIndex)

			err := p.processFunc(work)
			if err != nil {
				log.Printf("Pool %s: error processing work for worker %d: %v", p.poolLabel,
					workerIndex, err)
			}

			continue
		}

		// No tasks, sleep for a bit to avoid spinning.
		time.Sleep(time.Second)
	}

	log.Printf("Pool %s: stopped worker %d.", p.poolLabel, workerIndex)
}
