package mapreduce

type RegisterArgs struct {
	WorkerHostname string
}

type RegisterReply struct {
	WorkerId int
}

type HeartBeatArgs struct {
	WorkerId int
}

type HeartBeatReply struct {
	Ok bool
}

type RunMapArgs struct {
	MapId    int
	FilePath string
}

type RunReduceArgs struct {
	ReduceId int
	FilePath string
}

type MapDoneArgs struct {
	WorkerId int
	MapId    int
}

type ReduceDoneArgs struct {
	WorkerId int
	ReduceId int
}
