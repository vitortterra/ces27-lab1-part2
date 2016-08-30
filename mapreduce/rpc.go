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
	Beat int64
}

type HeartBeatReply struct {
	Ok bool
}
