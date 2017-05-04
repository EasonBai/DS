package main

import (
	"fmt"
	//"net/rpc"
	"net"
	"net/rpc"
	"net/http"
	//"log"
)

type Client struct{
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
	c := new(Client)

	args := C2SArgs{ClientId : 0, FileId : "test.txt", OpType: 101}

	reply := &C2SReply{}
	ok := c.sendOpToServer("127.0.0.1", "7005", args, reply)
	if ok{
		fmt.Println(reply.Success)
	}
}

func (client *Client) listening(){
	rpc.Register(client)
	rpc.HandleHTTP()


	l, err := net.Listen("tcp", ":"+client.localPort)
	if err != nil {
		fmt.Println("listen error")
	}
	for{
		http.Serve(l, nil)
	}

}

//For client, call server's Operation func
func (client *Client)sendOpToServer(serverIP string, serverPort string, args C2SArgs, reply *C2SReply) bool{
	server, err := rpc.DialHTTP("tcp", serverIP+":"+serverPort)
	if err != nil{
		fmt.Println("dialing error")
		return false
	}

	err = server.Call("Server.Operate", args,reply)

	if err != nil{
		fmt.Println("server error")
		fmt.Println(err)
		return false
	}

	return true
	
}