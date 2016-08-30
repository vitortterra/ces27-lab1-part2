package mapreduce

import (
	"log"
	"net"
	"net/rpc"
	"time"
)

type Worker struct {
	id             int
	hostname       string
	masterHostname string
	listener       net.Listener
	rpcServer      *rpc.Server
	task           *Task
}

func (worker *Worker) echo(msg string) {
	args := &EchoArgs{msg}

	err := worker.callMaster("Master.Echo", args, new(struct{}))

	if err != nil {
		log.Println("Echo failed. Error:", err)
	}
}

func (worker *Worker) register() error {
	var (
		err   error
		args  *RegisterArgs
		reply *RegisterReply
	)

	args = new(RegisterArgs)
	args.WorkerHostname = worker.hostname
	reply = new(RegisterReply)
	err = worker.callMaster("Master.Register", args, reply)

	worker.id = reply.WorkerId
	log.Println("Worker Id:", worker.id)

	return err
}

func (worker *Worker) heartMonitor(hb int) {
	var (
		err error
	)

	for {
		log.Println("Sending HeartBeat")

		args := &HeartBeatArgs{worker.id}
		var reply HeartBeatReply
		err = worker.callMaster("Master.HeartBeat", args, &reply)

		if err != nil {
			log.Println("HeartBeat failed. Error:", err)
		}

		if !reply.Ok {
			log.Fatal("Not recognized by Master.")
		}

		time.Sleep(time.Duration(hb) * time.Second)
	}
}

func (worker *Worker) callMaster(proc string, args interface{}, reply interface{}) error {
	var (
		err    error
		client *rpc.Client
	)

	client, err = rpc.Dial("tcp", worker.masterHostname)

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

func (worker *Worker) handleConnection(conn *net.Conn) error {
	log.Println("Serving RPCS to", (*conn).RemoteAddr())
	worker.rpcServer.ServeConn(*conn)
	log.Println("Closing connection to", (*conn).RemoteAddr())
	(*conn).Close()
	return nil
}
