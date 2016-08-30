package mapreduce

import (
	"log"
	"net"
	"net/rpc"
	"net/url"
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
	return
}

func (master *Master) Echo(args *EchoArgs, _ *struct{}) error {
	log.Println(args.Msg)
	return nil
}

func (master *Master) Register(args *RegisterArgs, reply *RegisterReply) error {
	log.Println("Registering worker", args.WorkerHostname)

	master.workersMutex.Lock()
	master.workerCounter++
	master.workers[master.workerCounter] = &RemoteWorker{args.WorkerHostname, WORKER_IDLE}
	master.workersMutex.Unlock()

	*reply = RegisterReply{master.workerCounter}
	return nil
}

var globalDeath int = 0

func (master *Master) HeartBeat(args *HeartBeatArgs, reply *HeartBeatReply) error {
	var (
		workerExists bool = false
	)

	master.workersMutex.Lock()
	if _, ok := master.workers[args.WorkerId]; ok {
		workerExists = true
	} else {
		log.Println("Unrecognized Client with Id", args.WorkerId)
	}
	master.workersMutex.Unlock()

	*reply = HeartBeatReply{workerExists}
	return nil
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
	log.Println("Running map operations")
	for url := range task.InputURLChan {
		for {
			started := false

			master.workersMutex.Lock()
			for i := 1; i <= master.workerCounter; i++ {
				if master.workers[i].status == WORKER_IDLE {
					master.workers[i].status = WORKER_RUNNING
					go runMap(master.workers[i], url)
					started = true
					break
				}
			}
			master.workersMutex.Unlock()

			if !started {
				log.Println("No worker available.")
				time.Sleep(time.Duration(3) * time.Second)
			}
		}
	}
}

func runMap(remoteWorker *RemoteWorker, url *url.URL) {
	log.Println("Running", url.String())
	args := new(RunMapArgs)
	args.RawUrl = url.String()
	err := remoteWorker.callRemoteWorker("Worker.RunMap", args, new(struct{}))

	if err != nil {
		log.Println("RunMap Failed. Error:", err)
	}
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
