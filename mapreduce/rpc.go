package mapreduce

type EchoArgs struct {
	Msg string
}

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
	RawUrl string
}
