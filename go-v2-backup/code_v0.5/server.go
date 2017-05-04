package hahaha

import (

	"sync"
	"time"
	"fmt"
	"net"
	"bytes"
	"net"
	"net/rpc"
	"net/http"
)



type Server struct {
	mu      sync.Mutex
	me      int
	isMaster bool
	//connection *TCPListener
	FileOpTable	 map[string]FileOpRecord
	ClientOpTable map[int]ClientOpRecord
	SeenClientRid map[int][]int
	// Your definitions here.
	localPort string

	ClientList []ClientEnd
}


func (svr *Server) update(args *PutAppendArgs) {
	//update FileOpTable and ClientOpTable
	//sendAppendEntries&appendEntries: send updated tables to backup
	//replytoclient=ok
	...
}

func (svr *Server) appendEntrytoClients(cid int64){
	//look up ClientOpTable
	//if data!=null, send data to client
	//if client reply ok, update FileOpTable and ClientOpTable
	//sendAppendEntries&appendEntries: send updated tables to backup
	...


	//rpc call client.Operate()
	//sendOpToClient(client.ip, client.port, args, reply)
}

func (svr *Server) sendAppendEntries(server int, args *AppendEntriesArgs, reply *AppendEntriesReply) bool{
	// send all the information updated to backup server
	...
	ok :=...
	if !ok{
		...
	}
}

func (svr *Server) appendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) {
	//backup server append information from master
	...
}


//
// servers[] contains the ports of the set of
// servers that will cooperate via Raft to
// form the fault-tolerant key/value service.
// me is the index of the current server in servers[].
// the k/v server should store snapshots with persister.SaveSnapshot(),
// and Raft should save its state (including log) with persister.SaveRaftState().
// the k/v server should snapshot when Raft's saved state exceeds maxraftstate bytes,
// in order to allow Raft to garbage-collect its log. if maxraftstate is -1,
// you don't need to snapshot.
// StartKVServer() must return quickly, so it should start goroutines
// for any long-running work.
//
func StartServer(clients []ClientEnd, me int, isMaster bool, port string) *Server {

	fmt.Println("start listening")
	
	//infinate loop 

	
	svr := new(Server)
	svr.me = me
	//set one server to be master and the other to be backup
	svr.isMaster=isMaster
	// Your initialization code here.
	svr.FileOpTable=make(map[string]FileOpRecord)
	svr.ClientOpTable=make(map[int64]ClientOpRecord)
	svr.localPort = port
	svr.ClientList = clients

	//start listening (for ever and ever)
	go svr.listening()


	return svr
}



// RPC related
func (svr *Server)Operate(args *C2SArgs, reply *C2SReply) error{
	// fmt.Println("Server.Operate")
	// fmt.Println(args.ClientId)
	// fmt.Println(args.FileId)
	// fmt.Println(args.OpType)


	switch(args.OpType){
		case "UPDATE":

			break
		case "DELETE":
	}

	reply.Success = "hahahahaha"

	return nil
}

func (svr *Server)listening(){
	rpc.Register(svr)
	rpc.HandleHTTP()


	l, err := net.Listen("tcp", ":"+svr.localPort)
	if err != nil {
		fmt.Println("listen error")
	}
	for{
		http.Serve(l, nil)
	}
}

func (svr *Server)sendOpToClient(clientIP string, clientPort string, args *S2CArgs, reply *S2CReply) bool{
	client, err := rpc.DialHTTP("tcp", clientIP+":"+clientPort)
	if err != nil{
		fmt.Println("dialing error")
		return false
	}

	err = client.Call("Client.Operate", args,reply)

	if err != nil{
		fmt.Println("RPC failed")
		fmt.Println(err)
		return false
	}

	return true
	
}