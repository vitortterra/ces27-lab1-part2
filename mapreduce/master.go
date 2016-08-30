package mapreduce

import (
	"log"
	"net"
	"net/rpc"
	"path/filepath"
	"sync"
	"time"
)

type Master struct {
	workersMutex  sync.Mutex
	rpcServer     *rpc.Server
	address       string
	listener      net.Listener
	done          chan bool
	workers       map[int]*RemoteWorker
	workerCounter int
	mapCounter    int
	reduceCounter int
}

type RemoteWorker struct {
	hostname string
	status   workerStatus
}

type workerStatus string

const (
	WORKER_IDLE    workerStatus = "idle"
	WORKER_RUNNING workerStatus = "running"
)

func newMaster(address string) (master *Master) {
	master = new(Master)
	master.address = address
	master.done = make(chan bool)
	master.workers = make(map[int]*RemoteWorker)
	master.workerCounter = 0
	master.mapCounter = 0
	master.reduceCounter = 0
	return
}

func (master *Master) acceptMultipleConnections() error {
	var (
		err     error
		newConn net.Conn
	)

	log.Println("Accepting connections on", master.listener.Addr())

	for {
		newConn, err = master.listener.Accept()

		if err == nil {
			go master.handleConnection(&newConn)
		} else {
			log.Println("Failed to accept connection. Error: ", err)
			break
		}
	}

	log.Println("Stopped accepting connections.")
	return nil
}

func (master *Master) handleConnection(conn *net.Conn) error {
	log.Println("Serving RPCS to", (*conn).RemoteAddr())
	master.rpcServer.ServeConn(*conn)
	log.Println("Closing connection to", (*conn).RemoteAddr())
	(*conn).Close()
	return nil
}

func (master *Master) runMaps(task *Task) {
	var wg sync.WaitGroup

	log.Println("Running map operations")
	for filePath := range task.InputFilePathChan {
		for {
			started := false

			master.workersMutex.Lock()
			for i := 1; i <= master.workerCounter; i++ {
				if master.workers[i].status == WORKER_IDLE {
					master.workers[i].status = WORKER_RUNNING
					wg.Add(1)
					go runMap(master.workers[i], filePath, master.mapCounter, &wg)
					master.mapCounter++
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
	log.Println("Map Completed")
}

func runMap(remoteWorker *RemoteWorker, filePath string, mapId int, wg *sync.WaitGroup) {
	log.Println("Running", filePath)
	args := &RunMapArgs{mapId, filePath}
	err := remoteWorker.callRemoteWorker("Worker.RunMap", args, new(struct{}))

	if err != nil {
		log.Println("RunMap Failed. Error:", err)
	}

	wg.Done()
	return
}

func (master *Master) runReduces(task *Task) {
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

func (worker *RemoteWorker) callRemoteWorker(proc string, args interface{}, reply interface{}) error {
	var (
		err    error
		client *rpc.Client
	)

	client, err = rpc.Dial("tcp", worker.hostname)

	if err != nil {
		return err
	}

	defer client.Close()

	err = client.Call(proc, args, reply)

	if err != nil {
		return err
	}

	return nil
}
