package mapreduce

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	RESULT_PATH = "result/"
)

// RPC - RunMap
// Run the map operation defined in the task and return when it's done.
func (worker *Worker) RunMap(args *RunMapArgs, _ *struct{}) error {
	var (
		err       error
		buffer    []byte
		mapResult []KeyValue
	)

	log.Printf("Running map id: %v, path: %v\n", args.MapId, args.FilePath)

	if buffer, err = ioutil.ReadFile(args.FilePath); err != nil {
		log.Fatal(err)
	}

	if worker.shouldFail(true) {
		panic("Induced failure.")
	}

	mapResult = worker.task.Map(buffer)
	storeLocal(worker.task, args.MapId, mapResult)
	return nil
}

// RPC - RunMap
// Run the reduce operation defined in the task and return when it's done.
func (worker *Worker) RunReduce(args *RunReduceArgs, _ *struct{}) error {
	log.Printf("Running reduce id: %v, path: %v\n", args.ReduceId, args.FilePath)

	var (
		err          error
		reduceResult []KeyValue
		file         *os.File
		fileEncoder  *json.Encoder
	)

	data := loadLocal(args.ReduceId)

	reduceResult = worker.task.Reduce(data)

	if worker.shouldFail(true) {
		panic("Induced failure.")
	}

	if file, err = os.Create(resultFileName(args.ReduceId)); err != nil {
		log.Fatal(err)
	}

	fileEncoder = json.NewEncoder(file)

	for _, value := range reduceResult {
		fileEncoder.Encode(value)
	}

	file.Close()
	return nil
}

// RPC - Done
// Will be called by Master when the task is done.
func (worker *Worker) Done(_ *struct{}, _ *struct{}) error {
	log.Println("Done.")
	defer func() {
		close(worker.done)
	}()
	return nil
}

// Support function to generate the name of result files
func resultFileName(id int) string {
	return filepath.Join(RESULT_PATH, fmt.Sprintf("result-%v", id))
}
