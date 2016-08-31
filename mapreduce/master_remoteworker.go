package mapreduce

import (
	"net/rpc"
)

type workerStatus string

const (
	WORKER_IDLE    workerStatus = "idle"
	WORKER_RUNNING workerStatus = "running"
)

type RemoteWorker struct {
	id       int
	hostname string
	status   workerStatus
}

// Call a RemoteWork with the procedure specified in parameters. It will also handle connecting
// to the server and closing it afterwards.
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
