package hahaha


import (
	"fmt"
	"strconv"
	"os"
	"bufio"
	"io/ioutil"
	"sync"
	"math/big"
	"crypto/rand"
	"net"
	"net/rpc"
	"net/http"
)
//import "fmt"


type Client struct {
	mu      sync.Mutex
	ServerList map[int]ServerEnd
	// You will have to modify this struct.
	id int
	masterId int
	ReqId int


}

func MakeClient(servers map[int]ServerEnd, me int, masterId int) *Client {
	clt := new(Client)
	clt.ServerList = servers
	// You'll have to add code here.
	clt.id=me
	clt.masterId=masterId
	clt.ReqId=1
	//clt.FileList=make(map[string]int)
	//ip := "127.0.0.1"
	//port := "7005"
	clt.localPort :=...
	/*
	* connect to server to indicate the start
	*/
	
	go clt.listening();
	go clt.start()

	return clt
}


func (clt *Client) start(){
	//receive and deal with requests from terminal
	for{
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Welcome to ....\n")
		fmt.Print("Please enter your command: \n")
		fmt.Print("1. create a new file \n")
		fmt.Print("2. update an existing file \n")
		fmt.Print("3. delete an existing file \n")
		St, _ := reader.ReadString('\n')
		St=strings.TrimSpace(St)
		St=St[:1]
		op,err:=strconv.Atoi(St)
		if err !=nil {
			fmt.Println("Wrong command. Try again.")
		} else{
			if op==1 {
				fmt.Println("Please enter file path")
				St, _ = reader.ReadString('\n')
				St=strings.TrimSpace(St)
				Stlist :=strings.Split(St, "\\")
				filename:=Stlist[len(Stlist)-1]
				go clt.update(St, filename, clt.ReqId)
				fmt.Println("Request %v received", clt.ReqId)
				clt.ReqId++
			} else if op==2{
				fmt.Println("Please enter file path")
				St, _ = reader.ReadString('\n')
				St=strings.TrimSpace(St)
				fmt.Println("Please enter the name of the file you want to replace")
				filename, _ := reader.ReadString('\n')
				filename=strings.TrimSpace(filename)
				go clt.update(St, filename, clt.ReqId)
				fmt.Println("Request %v received", clt.ReqId)
				clt.ReqId++
			} else if op==3{
				fmt.Println("Please enter filename")
				filename, _ := reader.ReadString('\n')
				filename=strings.TrimSpace(filename)
				go clt.delete(filename, clt.ReqId)
				fmt.Println("Request %v received", clt.ReqId)
				clt.ReqId++
			} else {
				fmt.Println("Wrong command. Try again.")
			}
		}
	}
}


func (clt *Client) delete(key string, rid int){

	// You will have to modify this function.
	//connectServer and send information
	...

}


func (clt *Client) update(key string, filename string, rid int) {
	// You will have to modify this function.
	//connectServer and send information
	err, content:=readFile(key)
	if err==false {
		fmt.Println("Request %v rejected due to wrong path", rId)
		return
	}
	args := &C2SMessage{}
	args.ClientId=clt.id
	args.ReqId=rid
	args.Filename=filename
	args.OpType=UPDATE
	args.Content=content
	//initialize all parameters before send
	reply := &C2SReply{}
	ok := clt.sendOpToServer(clt.ServerList[clt.masterId].IP, clt.ServerList[clt.masterId].localPort, args, reply)
	if ok{
		handleReply(rid)
		return
	}
	for !ok{
		select {
		case <-time.After(300 * time.Millisecond):
			ok= clt.sendOpToServer(clt.ServerList[clt.masterId].IP, clt.ServerList[clt.masterId].localPort, args, reply)
		case <-time.After(1200 * time.Millisecond):
			fmt.
		}
	}
}

func handleReply(rid int){

}

func readFile(path string) (bool, []byte){


	dat, err := ioutil.ReadFile(path)
	if err != nil{
		fmt.Println(err)
		return false, nil
	}

	return true, dat

}

func writeFile(path string, content []byte) bool{

	err := ioutil.WriteFile(path, content, 0644)

	if err != nil{
		fmt.Println(err)
		return false
	}

	return true
}

func (clt *Client)listening(){
	rpc.Register(clt)
	rpc.HandleHTTP()


	l, err := net.Listen("tcp", ":"+clt.localPort)
	if err != nil {
		fmt.Println("listen error")
	}
	for{
		http.Serve(l, nil)
	}
}

func (clt *Client)sendOpToServer(serverIP string, serverPort string, args *C2SArgs, reply *C2SReply) bool{
	server, err := rpc.DialHTTP("tcp", serverIP+":"+serverPort)
	if err != nil{
		fmt.Println("dialing error")
		return false
	}

	err = server.Call("Server.Operate", args,reply)

	if err != nil{
		fmt.Println("RPC failed")
		fmt.Println(err)
		return false
	}

	return true
	
}


func (clt *Client)Operate(args *S2CArgs, reply *S2CReply) error{


	switch(args.OpType){
		case "UPDATE":
			update(S2CArgs)
		case "DELETE":
	}

	//
	reply.Success = "hahahahaha"

	return nil
}