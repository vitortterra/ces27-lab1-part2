// common_rpc.go defined all the parameters used in RPC between
// master and workers
package mapreduce

type RegisterArgs struct {
	WorkerHostname string
}

type RegisterReply struct {
	WorkerId int
}

type RunMapArgs struct {
	MapId    int
	FilePath string
}

type RunReduceArgs struct {
	ReduceId int
	FilePath string
}
