package mapreduce

import (
	"log"
)

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

func (master *Master) MapDone(args *MapDoneArgs, _ *struct{}) error {
	log.Println("Map", args.MapId, "completed.")

	master.workersMutex.Lock()
	master.workers[args.WorkerId].status = WORKER_IDLE
	master.workersMutex.Unlock()

	return nil
}

func (master *Master) ReduceDone(args *ReduceDoneArgs, _ *struct{}) error {
	log.Println("Reduce", args.ReduceId, "completed.")

	master.workersMutex.Lock()
	master.workers[args.WorkerId].status = WORKER_IDLE
	master.workersMutex.Unlock()

	return nil
}

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
