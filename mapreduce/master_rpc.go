package mapreduce

import (
	"log"
)

// RPC - Register
// Procedure that will be called by workers to register within this master.
func (master *Master) Register(args *RegisterArgs, reply *RegisterReply) error {
	log.Printf("Registering worker %v with hostname %v", master.totalWorkers, args.WorkerHostname)

	master.workersMutex.Lock()
	master.workers = append(master.workers, &RemoteWorker{master.totalWorkers, args.WorkerHostname, WORKER_IDLE})
	master.workersMutex.Unlock()

	*reply = RegisterReply{master.totalWorkers}

	master.totalWorkers++
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

// RPC - HeartBeat
// Procedure that will be called by workers to send a heartbeat to master.
func (master *Master) HeartBeat(args *HeartBeatArgs, reply *HeartBeatReply) error {
	var (
		workerExists bool = false
	)

	master.workersMutex.Lock()
	if _, ok := master.workers[args.WorkerId]; ok {
		workerExists = true
	} else {
		log.Printf("Unrecognized Client with Id %v\n", args.WorkerId)
	}
	master.workersMutex.Unlock()

	*reply = HeartBeatReply{workerExists}
	return nil
}
