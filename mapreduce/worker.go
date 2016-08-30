package mapreduce

import (
	"log"
	"net"
	"net/rpc"
	"time"
)

type Worker struct {
	id int

	// Network
	hostname       string
	masterHostname string
	listener       net.Listener
	rpcServer      *rpc.Server

	// Operation
	task *Task
}

// Call RPC Register on Master to notify that this worker is ready
// to receive operations.
func (worker *Worker) register() error {
	var (
		err   error
		args  *RegisterArgs
		reply *RegisterReply
	)

	log.Println("Registering with Master")

	args = new(RegisterArgs)
	args.WorkerHostname = worker.hostname

	reply = new(RegisterReply)

	err = worker.callMaster("Master.Register", args, reply)

	worker.id = reply.WorkerId
	log.Printf("Registered. WorkerId: %v\n", worker.id)

	return err
}

// heartMonitor will repeatedly send information about this worker to Master.
func (worker *Worker) heartMonitor(hb int) {
	var (
		err   error
		args  *HeartBeatArgs
		reply *HeartBeatReply
	)

	for {
		log.Println("Sending HeartBeat")

		args = &HeartBeatArgs{worker.id}
		reply = new(HeartBeatReply)

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

// Handle a single connection until it's done, then closes it.
func (worker *Worker) handleConnection(conn *net.Conn) error {
	worker.rpcServer.ServeConn(*conn)
	(*conn).Close()
	return nil
}

// Connect to Master and call remote procedure.
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
