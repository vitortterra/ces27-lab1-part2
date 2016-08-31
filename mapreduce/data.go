package mapreduce

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	REDUCE_PATH = "reduce/"
)

// Returns the name of files created after merge
func mergeReduceName(idReduce int) string {
	return fmt.Sprintf("reduce-%v", idReduce)
}

// Returns the name of files created after map
func reduceName(idMap int, idReduce int) string {
	return fmt.Sprintf("reduce-%v-%v", idMap, idReduce)
}

// Store result from map operation locally.
// This will store the result from all the map calls.
func storeLocal(task *Task, idMapTask int, data []KeyValue) {
	var (
		err         error
		file        *os.File
		fileEncoder *json.Encoder
	)

	for r := 0; r < task.NumReduceJobs; r++ {
		file, err = os.Create(filepath.Join(REDUCE_PATH, reduceName(idMapTask, r)))
		if err != nil {
			log.Fatal(err)
		}

		fileEncoder = json.NewEncoder(file)
		for _, kv := range data {
			if task.Shuffle(task, kv.Key) == r {
				err = fileEncoder.Encode(&kv)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
		file.Close()
	}
}

// Merge the result from all the map operations by reduce job id.
func mergeLocal(task *Task, mapCounter int) {
	var (
		err              error
		file             *os.File
		fileDecoder      *json.Decoder
		mergeFile        *os.File
		mergeFileEncoder *json.Encoder
	)

	for r := 0; r < task.NumReduceJobs; r++ {
		if mergeFile, err = os.Create(filepath.Join(REDUCE_PATH, mergeReduceName(r))); err != nil {
			log.Fatal(err)
		}

		mergeFileEncoder = json.NewEncoder(mergeFile)

		for m := 0; m < mapCounter; m++ {
			if file, err = os.Open(filepath.Join(REDUCE_PATH, reduceName(m, r))); err != nil {
				log.Fatal(err)
			}

			fileDecoder = json.NewDecoder(file)

			for {
				var kv KeyValue
				err = fileDecoder.Decode(&kv)
				if err != nil {
					break
				}

				mergeFileEncoder.Encode(&kv)
			}
			file.Close()
		}

		mergeFile.Close()
	}
}

// Load data for reduce jobs.
func loadLocal(idReduce int) (data []KeyValue) {
	var (
		err         error
		file        *os.File
		fileDecoder *json.Decoder
	)

	if file, err = os.Open(filepath.Join(REDUCE_PATH, mergeReduceName(idReduce))); err != nil {
		log.Fatal(err)
	}

	fileDecoder = json.NewDecoder(file)

	data = make([]KeyValue, 0)

	for {
		var kv KeyValue
		if err = fileDecoder.Decode(&kv); err != nil {
			break
		}

		data = append(data, kv)
	}

	file.Close()
	return data
}
