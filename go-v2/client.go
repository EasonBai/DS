package main


import (
	"fmt"
	"strconv"
	"os"
	"bufio"
	"io/ioutil"
	"sync"
	//"math/big"
	//"crypto/rand"
	"net"
	"net/rpc"
	"net/http"
	"strings"
	"time"
	//"icommon"
)
//import "fmt"


type Client struct {
	mu      	sync.Mutex
	ServerList 	map[int]ServerEnd
	// You will have to modify this struct.
	id 			int
	masterId 	int
	ReqId 		int
	FileList 	map[string]int
	killed 		bool
	localPort 	string
}

func (clt *Client)MakeClient(servers map[int]ServerEnd, me int, masterId int, port string) {
	//clt := new(Client)
	clt.ServerList = servers
	// You'll have to add code here.
	clt.id=me
	clt.masterId=masterId
	clt.ReqId=1
	clt.FileList=make(map[string]int)
	clt.localPort=port
	clt.killed = false

	go clt.listening()
	clt.start()

	//return clt
}


func (clt *Client) start(){
	//receive and deal with requests from terminal
	clt.enter()
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Welcome to ....\n")
		fmt.Print("Please enter your command: \n")
		fmt.Print("1. create a new file \n")
		fmt.Print("2. update an existing file \n")
		fmt.Print("3. delete an existing file \n")
		fmt.Print("4. exit \n")
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
			} else if op==4 {
				clt.exit()
				clt.killed=true
				break
			} else {
				fmt.Println("Wrong command. Try again.")
			}
		}
	}
}


func (clt *Client) delete(key string, rid int){
	//connectServer and send information
	_, err := os.Stat("src\\"+ key)
	if os.IsNotExist(err) {
		fmt.Printf("Request %v rejected due to wrong filename \n", rid)
		return
	}
	args := &C2SArgs{}
	args.ClientId=clt.id
	args.ReqId=rid
	args.Filename=key
	args.OpType=DELETE
	reply := &C2SReply{}
	ok := clt.sendOpToServer(clt.ServerList[clt.masterId].IP, clt.ServerList[clt.masterId].LocalPort, args, reply)
	if ok{
		clt.handleReply(rid, reply)
		return
	}
	for !ok{
		select {
		case <-time.After(300 * time.Millisecond):
			ok= clt.sendOpToServer(clt.ServerList[clt.masterId].IP, clt.ServerList[clt.masterId].LocalPort, args, reply)
		case <-time.After(1200 * time.Millisecond):
			fmt.Println("Cannot connect to server. Check your net connection and try again")
			break
		}
	}
}


func (clt *Client) update(key string, filename string, rid int) {
	//connectServer and send information
	err, content:=readFile(key)
	if err==false {
		fmt.Println("Request %v rejected due to wrong path", rid)
		return
	}
	args := &C2SArgs{}
	args.ClientId=clt.id
	args.ReqId=rid
	args.Filename=filename
	args.OpType=UPDATE
	args.Content=content
	//initialize all parameters before send
	reply := &C2SReply{}
	ok := clt.sendOpToServer(clt.ServerList[clt.masterId].IP, clt.ServerList[clt.masterId].LocalPort, args, reply)
	if ok{
		clt.handleReply(rid, reply)
		return
	}
	for !ok{
		select {
		case <-time.After(300 * time.Millisecond):
			ok= clt.sendOpToServer(clt.ServerList[clt.masterId].IP, clt.ServerList[clt.masterId].LocalPort, args, reply)
		case <-time.After(1200 * time.Millisecond):
			fmt.Println("Cannot connect to server. Check your net connection and try again")
			break
		}
	}
}

func (clt *Client)handleReply(rid int, reply *C2SReply){
	if reply.Success{
		fmt.Printf("Request %v is accepted. \n", rid)
	} else if reply.Msg==WRONGLEAD {
		clt.masterId=(clt.masterId+1)%2
		fmt.Println("Request %v is rejected due to wrong leader.", rid)
	} else if reply.Msg==OLDREQ {
		fmt.Println("Request %v is out of date.", rid)
	} else {
		fmt.Println("Request %v is rejected."+reply.Msg, rid)
	}
}

func (clt *Client) enter() bool{
	args := C2SArgs{ClientId: clt.id}
	reply := &C2SReply{}
	return clt.sendEnterToServer(clt.ServerList[clt.masterId].IP,clt.ServerList[clt.masterId].LocalPort, args, reply) && reply.Success
}

func (clt *Client) sendEnterToServer(masterIP string, masterPort string, args C2SArgs, reply *C2SReply) bool{

	master, err := rpc.DialHTTP("tcp", masterIP+":"+masterPort)
	if err != nil{
		fmt.Println("dialing error")
		return false
	}

	err = master.Call("Server.Enter", args, reply)

	if err != nil{
		fmt.Println("server error")
		fmt.Println(err)
		return false
	}

	return true
}

func (clt *Client) exit() bool{
	args := C2SArgs{ClientId: clt.id}
	reply := &C2SReply{}
	return clt.sendExitToServer(clt.ServerList[clt.masterId].IP, clt.ServerList[clt.masterId].LocalPort, args, reply) && reply.Success

}
func (clt *Client) sendExitToServer(masterIP string, masterPort string, args C2SArgs, reply*C2SReply) bool{
	master, err := rpc.DialHTTP("tcp", masterIP+":"+masterPort)
		if err != nil{
		fmt.Println("dialing error")
		return false
	}

	err = master.Call("Server.Exit", args, reply)

	if err != nil{
		fmt.Println("server error")
		fmt.Println(err)
		return false
	}

	return true
}

func (clt *Client)HandleServerRequest(args *S2CArgs, reply *S2CReply){
	clt.mu.Lock()
	defer clt.mu.Unlock()
	fmt.Println("receive args from server")
	if args.TimeStamp <clt.FileList[args.FileId]{
		reply.Success=false
		reply.Msg=DUPLICATE
		return
	}
	switch args.OpType{
	case UPDATE:
		err := os.Mkdir("src", 0777)
		if err !=nil{
		}
		path:="src\\"+args.FileId
		check:=writeFile(path, args.Content)
		if check {
			clt.FileList[args.FileId]=args.TimeStamp
			reply.Success=true
			fmt.Println("Success update")
			reply.Msg=WRITTEN
		} else {
			reply.Success=false
			reply.Msg=UNABLEWRITE
		}
	case DELETE:
		err:=os.Mkdir("src", 0777)
		err=os.Remove("src\\"+args.FileId)
		if err==nil{
			clt.FileList[args.FileId]=args.TimeStamp
			reply.Success=true
			reply.Msg=REMOVED
			fmt.Println("Success removed")
		} else {
			clt.FileList[args.FileId]=args.TimeStamp
			reply.Success=true
			reply.Msg=REMOVED
		}
	}
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
	//rpc.HandleHTTP()
	rpc.HandleHTTP()

	l, err := net.Listen("tcp", ":"+clt.localPort)
	if err != nil {
		fmt.Println("listen error")
	}
	for !clt.killed {
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

func (clt *Client)CheckIfMasterRun(args *S2CArgs, reply *S2CReply){
	//fmt.Println("now check if master work!")
	cArgs := &C2SArgs{clt.id, 0, "", "CHECK", []byte("Check")}
	cReply := &C2SReply{}
	isOK := clt.sendOpToServer(clt.ServerList[clt.masterId].IP, clt.ServerList[clt.masterId].LocalPort, cArgs,cReply)
	if isOK == false {
		reply.Success = true
	}
}

func (clt *Client)Operate(args *S2CArgs, reply *S2CReply) error{
	switch args.OpType {
		case "CHECK":
			clt.CheckIfMasterRun(args, reply)
		default :
			clt.HandleServerRequest(args, reply)
	}
	return nil
}