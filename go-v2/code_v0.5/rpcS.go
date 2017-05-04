package main

import (
	"fmt"
	//"net/rpc"
	"net"
	"net/rpc"
	"net/http"
	//"log"
)

type Server struct{
	localPort string
}

type C2SArgs struct{
	ClientId 	int64
	FileId		string
	OpType		int
	Content 	[]byte
}

type C2SReply struct{
	Success string
}

type S2CArgs struct{
	ServerId	int64
	FileId		string
	OpType		int 
	TimeStamp 	int64
	Content 	[]byte
}

type S2CReply struct{
	Success string
}

func main(){

	s := new(Server)
	s.localPort = "7005"

	s.listening()
}

func (server *Server)listening(){
	rpc.Register(server)
	rpc.HandleHTTP()


	l, err := net.Listen("tcp", ":"+server.localPort)
	if err != nil {
		fmt.Println("listen error")
	}
	for{
		http.Serve(l, nil)
	}
}

//For server, call client's Operation func
// func (server *Server)sendOpToClient(clientIP string,clientPort string,  args S2CArgs, reply *S2CReply) bool {
	
// 	client, err := rpc.DialHTTP("tcp", clientIP + ":"+clientPort)
// 	if err != nil{
// 		fmt.Println("dialing error")
// 		return false
// 	}

// 	err = client.Call("Client.Operate", args,reply)

// 	if err != nil{
// 		fmt.Println("server error")
// 		return false
// 	}

// 	return true
// }

//RPC function
func (server *Server)Operate(args C2SArgs, reply *C2SReply) error{
	fmt.Println("Server.Operate")
	fmt.Println(args.ClientId)
	fmt.Println(args.FileId)
	fmt.Println(args.OpType)


	switch(args.OpType){
		
	}

	reply.Success = "hahahahaha"

	return nil
}