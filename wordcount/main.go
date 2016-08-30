package main

import (
	"flag"
	"github.com/pauloaguiar/ces27-lab1-part2/mapreduce"
	"log"
	"net/url"
	"os"
	"strconv"
)

var (
	// Run mode configuration settings
	mode       = flag.String("mode", "distributed", "Run mode: distributed or split")
	nodeType   = flag.String("type", "worker", "Node type: master or worker")
	reduceJobs = flag.Int("reducejobs", 5, "Number of reduce jobs that should be run")

	// Input data settings
	file      = flag.String("file", "files/pg1342.txt", "File to use as input")
	chunkSize = flag.Int("chunksize", 10*1024, "Size of data chunks that should be passed to map jobs(in bytes)")

	// Network settings
	addr   = flag.String("addr", "localhost", "IP address to listen on")
	port   = flag.Int("port", 5000, "TCP port to listen on")
	master = flag.String("master", "localhost:5000", "Master address")
)

// Entry point
func main() {
	var (
		err      error
		task     *mapreduce.Task
		numFiles int
		hostname string
	)

	flag.Parse()

	log.Println("Running in", *mode, "mode.")

	_ = os.Mkdir(MAP_PATH, os.ModeDir)
	_ = os.Mkdir(RESULT_PATH, os.ModeDir)

	switch *mode {
	case "distributed":
		switch *nodeType {
		case "master":
			var (
				fanIn chan *url.URL
			)
			log.Println("NodeType:", *nodeType)
			log.Println("Reduce Jobs:", *reduceJobs)
			log.Println("Address:", *addr)
			log.Println("Port:", *port)
			log.Println("File:", *file)
			log.Println("Chunk Size:", *chunkSize)

			hostname = *addr + ":" + strconv.Itoa(*port)

			// Splits data into chunks with size up to chunkSize
			if numFiles, err = splitData(*file, *chunkSize); err != nil {
				log.Fatal(err)
			}

			// Create fan in and out channels for mapreduce.Task
			fanIn = fanInUrl(numFiles, hostname)

			// Initialize mapreduce.Task object with the channels created above and functions
			// mapFunc, shufflerFunc and reduceFunc defined in wordcount.go
			task = &mapreduce.Task{
				Map:           mapFunc,
				Shuffle:       shuffleFunc,
				Reduce:        reduceFunc,
				NumReduceJobs: *reduceJobs,
				InputURLChan:  fanIn,
			}

			mapreduce.RunMaster(task, hostname)

		case "worker":
			log.Println("NodeType:", *nodeType)
			log.Println("Address:", *addr)
			log.Println("Port:", *port)
			log.Println("Master:", *master)

			hostname = *addr + ":" + strconv.Itoa(*port)
			mapreduce.RunWorker(task, hostname, *master)
		}
	}
}
