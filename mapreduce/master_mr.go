package mapreduce

import (
	"log"
	"path/filepath"
	"sync"
)

// Schedules map operations on remote workers. This will run until InputFilePathChan
func (master *Master) scheduleMaps(task *Task) {
	// is closed. If there is no worker available, it'll block.
	var (
		wg        sync.WaitGroup
		filePath  string
		worker    *RemoteWorker
		operation *Operation
	)

	log.Println("Running map operations")

	master.retryOperation = make(chan *Operation, RETRY_OPERATION_BUFFER)

	for filePath = range task.InputFilePathChan {
		operation = &Operation{"Worker.RunMap", master.mapCounter, filePath}
		master.mapCounter++

		worker = <-master.idleWorkerChan
		wg.Add(1)
		go master.runOperation(worker, operation, &wg)
	}

	go func() {
		for operation := range master.retryOperation {
			worker = <-master.idleWorkerChan
			log.Printf("Retrying map operation %v\n", operation.id)
			go master.runOperation(worker, operation, &wg)
		}
	}()

	wg.Wait()
	close(master.retryOperation)

	log.Println("Map Completed")
}

// Schedules reduce operations on remote workers. This will run until reduceFilePathChan
// is closed. If there is no worker available, it'll block.
func (master *Master) scheduleReduces(task *Task) {
	var (
		wg                 sync.WaitGroup
		filePath           string
		worker             *RemoteWorker
		operation          *Operation
		reduceFilePathChan chan string
	)

	log.Println("Running reduce operations")

	reduceFilePathChan = fanReduceFilePath(task.NumReduceJobs)
	master.retryOperation = make(chan *Operation, RETRY_OPERATION_BUFFER)

	for filePath = range reduceFilePathChan {
		operation = &Operation{"Worker.RunReduce", master.reduceCounter, filePath}
		master.reduceCounter++

		worker = <-master.idleWorkerChan
		wg.Add(1)
		go master.runOperation(worker, operation, &wg)
	}

	go func() {
		for operation := range master.retryOperation {
			worker = <-master.idleWorkerChan
			log.Printf("Retrying reduce operation %v\n", operation.id)
			go master.runOperation(worker, operation, &wg)
		}
	}()

	wg.Wait()
	close(master.retryOperation)

	log.Println("Reduce Completed")
}

// runOperation start a single operation on a RemoteWorker and wait for it to return or fail.
func (master *Master) runOperation(remoteWorker *RemoteWorker, operation *Operation, wg *sync.WaitGroup) {
	var (
		err  error
		args *RunArgs
	)

	log.Printf("Running %v '%v' file '%v' on worker '%v'\n", operation.proc, operation.id, operation.filePath, remoteWorker.id)

	args = &RunArgs{operation.id, operation.filePath}
	err = remoteWorker.callRemoteWorker(operation.proc, args, new(struct{}))

	if err != nil {
		log.Printf("Operation %v '%v' Failed. Error: %v\n", operation.proc, operation.id, err)
		master.retryOperation <- operation
		master.failedWorkerChan <- remoteWorker
	} else {
		wg.Done()
		master.idleWorkerChan <- remoteWorker
	}
}

// FanIn is a pattern that will return a channel in which the goroutines generated here will keep
// writing until the loop is done.
// This is used to generate the name of all the reduce files.
func fanReduceFilePath(numReduceJobs int) chan string {
	var (
		outputChan chan string
		filePath   string
	)

	outputChan = make(chan string)

	go func() {
		for i := 0; i < numReduceJobs; i++ {
			filePath = filepath.Join(REDUCE_PATH, mergeReduceName(i))

			outputChan <- filePath
		}

		close(outputChan)
	}()
	return outputChan
}
