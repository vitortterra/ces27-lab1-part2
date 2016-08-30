package mapreduce

import (
	"log"
	"net"
	"net/rpc"
	"os"
)

// RunSequential will ensure that map and reduce function runs in
// a single-core linearly. The Task is passed from the calling package
// and should contains the definitions for all the required functions
// and parameters.
// Notice that this implementation will store data locally. In the distributed
// version of mapreduce it's common to store the data in the same worker that computed
// it and just pass a reference to reduce jobs so they can go grab it.
func RunSequential(task *Task) {
	var (
		mapCounter int = 0
		mapResult  []KeyValue
	)

	log.Print("Running RunSequential...")

	_ = os.Mkdir(REDUCE_PATH, os.ModeDir)

	for v := range task.InputChan {
		mapResult = task.Map(v)
		storeLocal(task, mapCounter, mapResult)
		mapCounter++
	}

	mergeLocal(task, mapCounter)

	for r := 0; r < task.NumReduceJobs; r++ {
		data := loadLocal(r)
		task.OutputChan <- task.Reduce(data)
	}

	close(task.OutputChan)
	return
}

func RunMaster(task *Task, hostname string) {
	var (
		err      error
		master   *Master
		rpcs     *rpc.Server
		listener net.Listener
	)

	log.Println("Running Master on", hostname)
	master = newMaster(hostname)

	rpcs = rpc.NewServer()
	rpcs.Register(master)
	master.rpcServer = rpcs

	listener, err = net.Listen("tcp", master.address)

	if err != nil {
		log.Panic(err)
	}

	master.listener = listener

	go master.acceptMultipleConnections()

	master.runMaps(task)

	<-master.done
	return
}

func RunWorker(task *Task, hostname string, masterHostname string) {
	var (
		err      error
		worker   *Worker
		rpcs     *rpc.Server
		listener net.Listener
	)

	log.Println("Running Worker on", hostname)

	worker = new(Worker)
	worker.hostname = hostname
	worker.masterHostname = masterHostname

	rpcs = rpc.NewServer()
	rpcs.Register(worker)

	worker.rpcServer = rpcs

	listener, err = net.Listen("tcp", worker.hostname)

	if err != nil {
		log.Panic("Starting RPC listener failed. Error:", err)
	}

	worker.listener = listener
	defer worker.listener.Close()

	err = worker.register()

	if err != nil {
		log.Panic("Register RPC failed. Error:", err)
	}

	go worker.heartMonitor(5)

	for {
		conn, err := worker.listener.Accept()
		if err == nil {
			go worker.handleConnection(&conn)
		} else {
			break
		}
	}
}
