package mapreduce

import (
	"log"
	"path/filepath"
	"sync"
	"time"
)

// Schedules map operations on remote workers. This will run until InputFilePathChan
// is closed. If there is no worker available, it'll sleep and retry.
func (master *Master) scheduleMaps(task *Task) {
	var (
		wg        sync.WaitGroup
		filePath  string
		worker    *RemoteWorker
		available bool
	)

	log.Println("Running map operations")
	for filePath = range task.InputFilePathChan {
		for {
			worker, available = master.getIdleWorker()

			if available {
				wg.Add(1)
				go worker.runMap(filePath, master.mapCounter, &wg)
				break
			} else {
				time.Sleep(time.Duration(1) * time.Second)
			}
		}
	}

	wg.Wait()
	log.Println("Map Completed")
}

func (master *Master) getIdleWorker() (*RemoteWorker, bool) {
	// Locks workers struct and searches for an idle worker.
	master.workersMutex.Lock()
	defer master.workersMutex.Unlock()
	for _, worker := range master.workers {
		if worker.status == WORKER_IDLE {
			return worker, true
		}
	}

	return nil, false
}

// runMap start a single map operation on a RemoteWorker and wait for it to return or fail.
func (remoteWorker *RemoteWorker) runMap(filePath string, mapId int, wg *sync.WaitGroup) {
	var (
		err  error
		args *RunMapArgs
	)

	log.Printf("Running Map '%v' file '%v' on worker '%v'\n", mapId, filePath, remoteWorker.id)

	args = &RunMapArgs{mapId, filePath}
	err = remoteWorker.callRemoteWorker("Worker.RunMap", args, new(struct{}))

	if err != nil {
		log.Println("Worker.RunMap Failed. Error:", err)
	}

	wg.Done()
}

func (master *Master) scheduleReduces(task *Task) {
	var wg sync.WaitGroup

	log.Println("Running reduce operations")
	for r := 0; r < task.NumReduceJobs; r++ {
		filePath := filepath.Join(REDUCE_PATH, mergeReduceName(r))

		for {
			started := false

			master.workersMutex.Lock()
			for i := 1; i <= master.workerCounter; i++ {
				if master.workers[i].status == WORKER_IDLE {
					master.workers[i].status = WORKER_RUNNING
					wg.Add(1)
					go runReduce(master.workers[i], filePath, master.reduceCounter, &wg)
					master.reduceCounter++
					started = true
					break
				}
			}
			master.workersMutex.Unlock()

			if !started {
				log.Println("No worker available.")
				time.Sleep(time.Duration(3) * time.Second)
			} else {
				break
			}
		}
	}

	wg.Wait()
	log.Println("Reduce Completed")
}

func runReduce(remoteWorker *RemoteWorker, filePath string, reduceId int, wg *sync.WaitGroup) {
	log.Println("Running", filePath)
	args := &RunReduceArgs{reduceId, filePath}
	err := remoteWorker.callRemoteWorker("Worker.RunReduce", args, new(struct{}))

	if err != nil {
		log.Println("RunReduce Failed. Error:", err)
	}

	wg.Done()
	return
}
