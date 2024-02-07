package kvsrv

import (
	"6.5840/labrpc"
	"sync"
	"sync/atomic"
)
import "crypto/rand"
import "math/big"

type Clerk struct {
	ID     int64
	server *labrpc.ClientEnd
	syn    int64
	mu     sync.Mutex
}

func nrand() int64 {
	max := big.NewInt(int64(1) << 62)
	bigx, _ := rand.Int(rand.Reader, max)
	x := bigx.Int64()
	return x
}

func MakeClerk(server *labrpc.ClientEnd) *Clerk {
	ck := new(Clerk)
	ck.ID = nrand()
	ck.server = server
	ck.syn = -1
	return ck
}

// fetch the current value for a key.
// returns "" if the key does not exist.
// keeps trying forever in the face of all other errors.
//
// you can send an RPC with code like this:
// ok := ck.server.Call("KVServer.Get", &args, &reply)
//
// the types of args and reply (including whether they are pointers)
// must match the declared types of the RPC handler function's
// arguments. and reply must be passed as a pointer.
func (ck *Clerk) Get(key string) string {
	// ck.mu.Lock()
	// defer ck.mu.Unlock()
	args := GetArgs{ck.ID, atomic.AddInt64(&ck.syn, 1), key}
	reply := GetReply{}
	for {
		ok := ck.server.Call("KVServer.Get", &args, &reply)
		if ok {
			// log.Printf("[Clerk] %d ReqNo: %d Get value from KVServer, key: %s, value: %s\n", ck.ID, args.ReqNo, key, reply.Value)
			return reply.Value
		}
		// log.Printf("[Clerk] Get: Unable to call KVServer, tried %d time(s)...\n", i+1)
	}
}

// shared by Put and Append.
//
// you can send an RPC with code like this:
// ok := ck.server.Call("KVServer."+op, &args, &reply)
//
// the types of args and reply (including whether they are pointers)
// must match the declared types of the RPC handler function's
// arguments. and reply must be passed as a pointer.
func (ck *Clerk) PutAppend(key string, value string, op string) string {
	// ck.mu.Lock()
	// defer ck.mu.Unlock()
	args := PutAppendArgs{ck.ID, atomic.AddInt64(&ck.syn, 1), key, value}
	reply := PutAppendReply{}
	for {
		ok := ck.server.Call("KVServer."+op, &args, &reply)
		if ok {
			// log.Printf("[Clerk] %d ReqNo: %d %s KV to KVServer successfully, key: %s, old value: %s new value: %s\n", ck.ID, args.ReqNo, op, key, reply.Value, reply.Value+value)
			return reply.Value
		}
		// log.Printf("[Clerk] %s: Unable to call KVServer, tried %d time(s)...\n", op, i+1)
	}
}

func (ck *Clerk) Put(key string, value string) {
	ck.PutAppend(key, value, "Put")
}

// Append value to key's value and return that value
func (ck *Clerk) Append(key string, value string) string {
	return ck.PutAppend(key, value, "Append")
}
