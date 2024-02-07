package kvsrv

// Put or Append
type PutAppendArgs struct {
	ID    int64
	ReqNo int64
	Key   string
	Value string
}

type PutAppendReply struct {
	Value string
}

type GetArgs struct {
	ID    int64
	ReqNo int64
	Key   string
}

type GetReply struct {
	Value string
}
