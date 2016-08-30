package mapreduce

import (
	"log"
	"net"
	"net/rpc"
	"sync"
	"time"
)

type Master struct {
	workersMutex  sync.Mutex
	rpcServer     *rpc.Server
	address       string
	listener      net.Listener
	done          chan bool
	workers       []RemoteWorker
	workerCounter int
}

type RemoteWorker struct {
	hostname string
}

func newMaster(address string) (master *Master) {
	master = new(Master)
	master.address = address
	master.done = make(chan bool)
	master.workers = make([]RemoteWorker, 0)
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
	master.workers = append(master.workers, RemoteWorker{args.WorkerHostname})
	master.workerCounter++
	master.workersMutex.Unlock()

	*reply = RegisterReply{master.workerCounter}
	return nil
}

func (master *Master) HeartBeat(args *HeartBeatArgs, reply *HeartBeatReply) error {
	*reply = HeartBeatReply{true}
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

func (master *Master) heartMonitor(hb int) {
	var (
		counter int64 = 0
		alive   int
		total   int
	)

	for {
		log.Println("Sending HeartBeat")
		alive = 0
		total = len(master.workers)
		for _, w := range master.workers {
			client, err := rpc.Dial("tcp", w.hostname)

			if err != nil {
				log.Println("Client dial failed. Err: ", err)
				continue
			}

			args := &HeartBeatArgs{counter}
			var reply HeartBeatReply
			err = client.Call("Worker.HeartBeat", args, &reply)

			if err != nil {
				log.Println("HeartBeat failed. Error:", err)
				continue
			}

			err = client.Close()
			if err != nil {
				log.Println("Close failed. Error:", err)
				continue
			}

			alive++
		}
		counter++
		log.Printf("Alive/Total: %v/%v\n", alive, total)
		time.Sleep(time.Duration(hb) * time.Second)
	}
}

func (master *Master) runMaps(task *Task) {
	for v := range task.InputChan {
		log.Println("Should schedule with len: ", len(v))
	}
}
