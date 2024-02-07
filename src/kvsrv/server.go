package kvsrv

import (
	"log"
	"sync"
)

const Debug = false

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug {
		log.Printf(format, a...)
	}
	return
}

type KVServer struct {
	mu           sync.Mutex
	mp           map[string]string
	duplicateRPC map[int64]map[int64]string
}

func (kv *KVServer) Get(args *GetArgs, reply *GetReply) {
	kv.mu.Lock()
	reply.Value = kv.mp[args.Key]
	//if _, exist := kv.duplicateRPC[args.ID][args.ReqNo]; !exist {
	//	reply.Value = kv.mp[args.Key]
	//	// reply.Value = kv.mp[args.Key]
	//	//if args.ReqNo == kv.expectReqNo[args.ID] {
	//	//	kv.expectReqNo[args.ID]++
	//	//	kv.mp[args.Key] = args.Value
	//	//	// reply.Value = kv.mp[args.Key]
	//	kv.addData(args.ID, args.ReqNo, reply.Value)
	//} else {
	//	reply.Value = kv.duplicateRPC[args.ID][args.ReqNo]
	//}
	//if args.ReqNo == kv.expectReqNo[args.ID] {
	//	kv.expectReqNo[args.ID]++
	//	reply.Value = kv.mp[args.Key]
	// kv.addData(args.ID, args.ReqNo, reply.Value)
	//	//kv.duplicateRPC[args.ID] = make(map[int64]string)
	//	//kv.duplicateRPC[args.ID][args.ReqNo] = reply.Value
	//} else if args.ReqNo < kv.expectReqNo[args.ID] {
	//	// log.Printf("[Server] Already comes, do nothing just reply")
	//	// reply.Value = kv.duplicateRPC[args.ID][args.ReqNo]
	//} else {
	//	// log.Printf("[Server] Not your turn, wait please ")
	//	reply.Wait = true
	//}
	// log.Println("[Server] Get", args.ID, "key:", args.Key, "value:", reply.Value, "ReqNo:", args.ReqNo)
	kv.mu.Unlock()
}

func (kv *KVServer) Put(args *PutAppendArgs, reply *PutAppendReply) {
	// log.Println("[Server] Put", args.ID, "key:", args.Key, "value:", args.Value, "ReqNo:", args.ReqNo)
	kv.mu.Lock()
	if _, exist := kv.duplicateRPC[args.ID][args.ReqNo]; !exist {
		kv.mp[args.Key] = args.Value
		// reply.Value = kv.mp[args.Key]
		//if args.ReqNo == kv.expectReqNo[args.ID] {
		//	kv.expectReqNo[args.ID]++
		//	kv.mp[args.Key] = args.Value
		//	// reply.Value = kv.mp[args.Key]
		kv.addData(args.ID, args.ReqNo, reply.Value)
	}
	//	//kv.duplicateRPC[args.ID] = make(map[int64]string)
	//	//kv.duplicateRPC[args.ID][args.ReqNo] = reply.Value
	//} else if args.ReqNo < kv.expectReqNo[args.ID] {
	//	// log.Printf("[Server] Already comes, do not thing just ")
	//	// reply.Value = kv.duplicateRPC[args.ID][args.ReqNo]
	//} else {
	//	// log.Printf("[Server] Not your turn, wait please ")
	//	reply.Wait = true
	//}
	kv.mu.Unlock()
}

func (kv *KVServer) Append(args *PutAppendArgs, reply *PutAppendReply) {
	// log.Println("[Server] Append", args.ID, "key:", args.Key, "value:", args.Value, "ReqNo:", args.ReqNo)
	kv.mu.Lock()
	if _, exist := kv.duplicateRPC[args.ID][args.ReqNo]; !exist {
		reply.Value = kv.mp[args.Key]
		kv.mp[args.Key] += args.Value
		kv.addData(args.ID, args.ReqNo, reply.Value)
	} else {
		reply.Value = kv.duplicateRPC[args.ID][args.ReqNo]
	}
	//if args.ReqNo == kv.expectReqNo[args.ID] {
	//	kv.expectReqNo[args.ID]++
	//	reply.Value = kv.mp[args.Key]
	//	kv.mp[args.Key] += args.Value
	//	kv.addData(args.ID, args.ReqNo, reply.Value)
	//	//kv.duplicateRPC[args.ID] = make(map[int64]string)
	//	//kv.duplicateRPC[args.ID][args.ReqNo] = reply.Value
	//} else if args.ReqNo < kv.expectReqNo[args.ID] {
	//	// log.Printf("[Server] Already comes, do not thing just ")
	//	reply.Value = kv.duplicateRPC[args.ID][args.ReqNo]
	//} else {
	//	// log.Printf("[Server] Not your turn, wait please ")
	//	reply.Wait = true
	//}
	kv.mu.Unlock()
}

func (kv *KVServer) addData(id int64, key int64, value string) {
	// 检查当前键对应的值数量
	if _, exists := kv.duplicateRPC[id]; !exists {
		// 如果键不存在，初始化内部的 map
		kv.duplicateRPC[id] = make(map[int64]string)
	}

	// 添加新的键值对
	kv.duplicateRPC[id][key] = value

	// 检查是否超过 20 个值
	if len(kv.duplicateRPC[id]) > 20 {
		// 删除最早添加的 10 个值
		// fmt.Println("deleting...")
		deleteOldValues(kv.duplicateRPC[id], key)
	}
}

// 删除最早添加的 10 个值
func deleteOldValues(m map[int64]string, curReqNo int64) {
	// 删除最早添加的 10 个值
	for i := curReqNo - 20; i < curReqNo-10; i++ {
		delete(m, i)
		// log.Printf("delete %d\n", i)
	}
}

func StartKVServer() *KVServer {
	kv := new(KVServer)
	kv.mp = make(map[string]string)
	kv.duplicateRPC = make(map[int64]map[int64]string)
	return kv
}
