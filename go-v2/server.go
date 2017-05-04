package main

import (

	"sync"
	"time"
	"fmt"
	"net"
	"strconv"
	//"bytes"
	"net/rpc"
	"net/http"
	"strings"
)



type Server struct {
	mu      sync.Mutex
	me      int
	isMaster bool
	FileOpTable	 map[string]FileOpRecord
	ClientOpTable map[int](map[string]struct{})
	SeenClientRid map[string]SeenRidRecord

	localPort string
	timestamp int
	peer	PeerAddress
	ClientList map[int]ClientEnd
	OnlineClientList map[int]struct{}
	AllClientsId map[int]struct{}
	killed bool
	masterHeartBeat chan bool
	ClientNumNotConMstr int
	SleepDuration 	int
}


func (svr *Server) update(args *C2SArgs, reply *C2SReply) {
	//update FileOpTable and ClientOpTable
	//sendAppendEntries&appendEntries: send updated tables to backup
	//replytoclient=ok
	filename:=args.Filename
	cid:=args.ClientId
	Rid:=args.ReqId
	key:=filename+"//"+strconv.Itoa(cid)
	fmt.Println(key)
	if svr.SeenClientRid[key].Rid>=Rid {
		reply.Msg=OLDREQ
		return
	}
	Entry := M2BArgs{}
	Entry.NewSeenRidRecord=SeenRidRecord{FilenameClient: key, Rid: Rid}
	if args.OpType==DELETE && svr.FileOpTable[filename].Status==DELETE{
		Entry.UpdateType=UPDATETYPE1		
	} else {
		svr.mu.Lock()
		Entry.UpdateType=UPDATETYPE2
		list:=make(map[int]struct{})
		for key := range svr.AllClientsId{
			list[key]=struct{}{}
		}
		Entry.NewFileOpRecord=FileOpRecord{Filename:filename, TimeStamp:svr.timestamp,Status:args.OpType,ClientIdList:list,Content:args.Content}
		svr.timestamp++
		svr.mu.Unlock()
	}
	BReply:=&M2BReply{}
	ok:=svr.sendRecordToBackup(svr.peer.IP, svr.peer.Port, Entry,BReply)
	if !ok || !BReply.Success{
		reply.Success=false
		reply.Msg=BACKUPLOST
		return
	}
	svr.mu.Lock()
	defer svr.mu.Unlock()
	switch Entry.UpdateType {
    case UPDATETYPE1:
    	if svr.SeenClientRid[key].Rid<Rid{
    		svr.SeenClientRid[key]=Entry.NewSeenRidRecord
    		reply.Success=true
    		reply.Msg=ACCEPTED
    	} else {
    		reply.Success=false
    		reply.Msg=OLDREQ
    	}
    case UPDATETYPE2:
    	if BReply.TimeStamp<=svr.FileOpTable[filename].TimeStamp{
    		if svr.SeenClientRid[key].Rid<Rid{
    			svr.SeenClientRid[key]=Entry.NewSeenRidRecord
    		}
    		reply.Success=false
    		reply.Msg=NOTSUCCESS
    	} else {
    		svr.SeenClientRid[key]=Entry.NewSeenRidRecord
    		svr.FileOpTable[filename]=Entry.NewFileOpRecord
    		fmt.Printf("num of all clients: %v \n",len(svr.AllClientsId))
    		for i:= range svr.AllClientsId {
    			fmt.Println("check if contains")
    			if _,contains :=svr.ClientOpTable[i][filename]; !contains{
    				svr.ClientOpTable[i][filename]= struct{}{}
    				fmt.Printf("added-id: %v, filename: %v \n", i, filename)
    			} else {
    				fmt.Printf("id %v already contains %v \n", i, svr.ClientOpTable[i][filename])
    			}
    		}
    		fmt.Printf("Client 0 Table length: %v \n ", len(svr.ClientOpTable[0]))
    		fmt.Printf("Client 1 Table length: %v \n ", len(svr.ClientOpTable[1]))
    		reply.Success=true
    		reply.Msg=ACCEPTED    		
    	}
    }
    svr.printTables()

}

func (svr *Server) sendEntrytoClients(){
	//look up ClientOpTable
	//if data!=null, send data to client
	//if client reply ok, update FileOpTable and ClientOpTable
	for id:=range svr.OnlineClientList{
		if len(svr.ClientOpTable[id])>0 {
			//go func(){
				IP:=svr.ClientList[id].IP
				Port:=svr.ClientList[id].LocalPort
				for filename:=range svr.ClientOpTable[id]{
					args:=&S2CArgs{}
					args.ClientId=id
					args.FileId=filename
					args.OpType=svr.FileOpTable[filename].Status
					args.TimeStamp=svr.FileOpTable[filename].TimeStamp
					args.Content=svr.FileOpTable[filename].Content
					reply := &S2CReply{}
					fmt.Println("send args to client")
					ok:=svr.sendOpToClient(IP,Port,args,reply)
					if ok{
						//go svr.handleClientsReply(args, reply)
						svr.handleClientsReply(args, reply)
					}
				}
			//}()
		}
	}

}

func (svr *Server) handleClientsReply(args *S2CArgs, reply *S2CReply){
	svr.mu.Lock()
	defer svr.mu.Unlock()
	if args.TimeStamp<svr.FileOpTable[args.FileId].TimeStamp{
		return
	}
	if reply.Success{
		Entry := M2BArgs{}
		Entry.UpdateType=3
		delete(svr.ClientOpTable[args.ClientId],args.FileId)
		Entry.NewClientOpRecord=ClientOpRecord{ClientId: args.ClientId, FileIdList: svr.ClientOpTable[args.ClientId]}
		delete(svr.FileOpTable[args.FileId].ClientIdList, args.ClientId)
		if len(svr.FileOpTable[args.FileId].ClientIdList)==0 {
			Entry.UpdateType=4
			Entry.NewFileOpRecord=FileOpRecord{Filename: args.FileId, TimeStamp: svr.FileOpTable[args.FileId].TimeStamp}
			delete(svr.FileOpTable, args.FileId)
		} else {
			Entry.NewFileOpRecord=svr.FileOpTable[args.FileId]
		}
		BReply:=&M2BReply{}
		svr.printTables()
		svr.sendRecordToBackup(svr.peer.IP, svr.peer.Port, Entry,BReply)
	}

}

func (svr *Server)sendRecordToBackup(backupIP string, backupPort string, args M2BArgs, reply *M2BReply) bool{
	backup, err := rpc.DialHTTP("tcp", backupIP+":"+backupPort)
	if err != nil{
		fmt.Println("sendRecordToBackup")
		fmt.Println("dialing error")
		return false
	}

	err = backup.Call("Server.BackupTables", args, reply)

	if err != nil{
		fmt.Println("server error")
		fmt.Println(err)
		return false
	}

	return true
}

func (svr *Server) BackupTables(args M2BArgs, reply *M2BReply) error{
	//fmt.Println(args.NewFileOpRecord.Filename)
	svr.mu.Lock()
	defer svr.mu.Unlock()



	if svr.isMaster {
		reply.Success = false
		return nil
	}
	
	go func(){svr.masterHeartBeat <- true}()
	
	
	reply.Success = true //default
	//newFileOpRecord = &args.NewFileOpRecord
	switch args.UpdateType{
	// case 1 finished!
	case 1:
		fmt.Println("case 1")
		//update seenclientrid table
		if args.NewSeenRidRecord.Rid > svr.SeenClientRid[args.NewSeenRidRecord.FilenameClient].Rid{
			svr.SeenClientRid[args.NewSeenRidRecord.FilenameClient] = args.NewSeenRidRecord
			//fmt.Printf("update rid table: %v\n",svr.SeenClientRid[args.NewSeenRidRecord.FilenameClient])
		}else{
			reply.Success = false
		}
	// case 2 finished!
	case 2:
		fmt.Println("case 2")
		reply.TimeStamp = args.NewFileOpRecord.TimeStamp
		// Not contained or smaller timestamp
		if oldRecord, contains := svr.FileOpTable[args.NewFileOpRecord.Filename]; !contains || oldRecord.TimeStamp < args.NewFileOpRecord.TimeStamp{
			//update seenclientrid table
			svr.SeenClientRid[args.NewSeenRidRecord.FilenameClient] = args.NewSeenRidRecord
			//fmt.Printf("update rid table: %v\n",svr.SeenClientRid[args.NewSeenRidRecord.FilenameClient])
			//update fileoptable
			svr.FileOpTable[args.NewFileOpRecord.Filename] = args.NewFileOpRecord
			//fmt.Printf("update file op table: %v\n", len(svr.FileOpTable))
			//update clientoptable
			for clientId,_ := range svr.ClientOpTable{
				svr.ClientOpTable[clientId][args.NewFileOpRecord.Filename] = struct{}{}
				//fmt.Printf("update client op table %v\n", len(svr.ClientOpTable[clientId]))
			}
		}else{
			reply.Success = false
		}
	// case 3 finished!	
	case 3:
		fmt.Println("case 3")
		reply.TimeStamp = args.NewFileOpRecord.TimeStamp
		
		//Contains && bigger timestamp
		if oldRecord,contains := svr.FileOpTable[args.NewFileOpRecord.Filename]; contains && oldRecord.TimeStamp <= args.NewFileOpRecord.TimeStamp{
			//update file op table
			newRecord := args.NewFileOpRecord
			if args.NewFileOpRecord.Status == UPDATE{
				newRecord.Content = svr.FileOpTable[newRecord.Filename].Content
			}
			svr.FileOpTable[newRecord.Filename] = newRecord

			//update clientoptable
			svr.ClientOpTable[args.NewClientOpRecord.ClientId] = args.NewClientOpRecord.FileIdList
			//fmt.Println("updated")
			//fmt.Printf("update client op table %v\n", len(svr.ClientOpTable[args.NewClientOpRecord.ClientId]))
		}else{
			reply.Success = false
		}
	// case 4 finished!
	case 4:
		fmt.Println("case 4")
		reply.TimeStamp = args.NewFileOpRecord.TimeStamp

		if oldRecord,contains := svr.FileOpTable[args.NewFileOpRecord.Filename]; contains && oldRecord.TimeStamp <= args.NewFileOpRecord.TimeStamp {
			delete(svr.FileOpTable, args.NewFileOpRecord.Filename)
			fmt.Println(len(svr.FileOpTable))
			svr.ClientOpTable[args.NewClientOpRecord.ClientId] = args.NewClientOpRecord.FileIdList
			//fmt.Printf("delete a record in Fileop table %v\n", len(svr.FileOpTable))
		}else{
			reply.Success = false
		}
	case 5:
		fmt.Println("case 5")
		//delete client from onlineclientlist
		delete(svr.OnlineClientList, args.OfflineClientId)
		fmt.Printf("OnlineClientList length: %v\n", len(svr.OnlineClientList))
		//reset rid to 0 in seen table
		for _, filenameclientid := range args.NewSeenRidRecordList{
			svr.SeenClientRid[filenameclientid] = SeenRidRecord{FilenameClient: filenameclientid, Rid: 0}
		}

		fmt.Printf("Rid with this client id: %v\n",svr.SeenClientRid["test.txt12//23"])		
	case 6:
		fmt.Println("case 6")

		svr.OnlineClientList[args.ClientId] = struct{}{}
	}

	svr.printTables()
	return nil
}

func main(me int, isMaster bool, port string, address PeerAddress) {
	
	fmt.Println("start listening")
	//svr := new(Server)
	svr.me = 0
	//set one server to be master and the other to be backup
	svr.isMaster = true
	svr.FileOpTable=make(map[string]FileOpRecord)
	svr.ClientOpTable=make(map[int](map[string]struct{}))
	svr.OnlineClientList=make(map[int]struct{})
	svr.AllClientsId=make(map[int]struct{})
	svr.SeenClientRid=make(map[string]SeenRidRecord)
	//for key := range clients{
		//svr.ClientOpTable[key]=make(map[string]struct{})
		//svr.AllClientsId[key]=struct{}{}
	//}
	svr.localPort = pp
	svr.ClientList = make(map[int]ClientEnd)
	svr.timestamp=1
	svr.peer=address
	svr.killed=false
	svr.ClientNumNotConMstr=0
	svr.SleepDuration=0
	svr.masterHeartBeat = make(chan bool)

	go svr.listening()
	svr.doStaff()

	//return svr
}

func(svr *Server)killing(){
	svr.killed = true
}

func (svr *Server) SetTestCase(args TestCaseArgs, reply *TestCaseReply) error{
	svr.mu.Lock()
	defer svr.mu.Unlock()
	switch args.Case {
	case 1:
		svr.killed = true
		fmt.Println("TEST CASE 1")
	case 2:
		svr.SleepDuration = args.SleepDuration
		fmt.Printf("TEST CASE 2 SleepDuration %v\n", args.SleepDuration)
	}

	return nil

}

func (svr *Server) Enter (args C2SArgs, reply *C2SReply) error{
	svr.mu.Lock()
	defer svr.mu.Unlock()

	fmt.Println(args.ClientId)

	//send update to backup
	m2bArgs := M2BArgs{UpdateType: 6, ClientId: args.ClientId}
	m2bReply := &M2BReply{}

	svr.sendRecordToBackup(svr.peer.IP, svr.peer.Port, m2bArgs, m2bReply)

	//add the client to the onlineclientlist

	svr.OnlineClientList[args.ClientId] = struct{}{}
	fmt.Println("onlineclientlist updated!")


	return nil
}

func (svr *Server)doStaff(){
	for !svr.killed{
		if svr.isMaster {
			svr.sendEntrytoClients()
		} else {
			svr.checkConnectionWithMaster()
		}
	}
}

func (svr *Server) checkConnectionWithMaster(){
	timeoutDuration := 5000
	select{
		case <-svr.masterHeartBeat:
		case <-time.After(time.Duration(timeoutDuration) * time.Millisecond):
			
			//timeout check master with all clients
			if (!svr.masterIsAlive()) && svr.Check(){
				svr.mu.Lock()
				svr.isMaster = true
				//fmt.Println("Backup becomes Master")
				svr.mu.Unlock()
			}
	}

}
func (svr *Server) CheckMasterIsAlive(args C2SArgs, reply *C2SReply) error{
	reply.Success = true

	return nil
}


func (svr *Server) masterIsAlive() bool {
	args := C2SArgs{}
	reply := &C2SReply{}
	reply.Success = false

	//TODO!!!!!!! change to master ip & port
	return svr.sendCheckMasterIsAlive(svr.peer.IP,svr.peer.Port,args,reply) && reply.Success
}

func (svr *Server) sendCheckMasterIsAlive(masterIP string, masterPort string, args C2SArgs, reply*C2SReply) bool{
	master, err := rpc.DialHTTP("tcp", masterIP+":"+masterPort)
		if err != nil{
		//fmt.Println("dialing error")
		return false
	}

	err = master.Call("Server.CheckMasterIsAlive", args, reply)

	if err != nil{
		//fmt.Println("server error")
		fmt.Println(err)
		return false
	}

	return true

}


func (svr *Server)CheckIfMasterDown(cend ClientEnd) {
	args := &S2CArgs{ ClientId: cend.Id , FileId: "", OpType: "CHECK", TimeStamp: 0,/*need to update*/ Content: []byte("Check")}
	reply := &S2CReply{}
	fmt.Println(cend.IP)
	fmt.Println(cend.LocalPort)
	isOK := svr.sendOpToClient(cend.IP, cend.LocalPort, args, reply)
	if isOK == false {
		fmt.Println("Failed to connect to client")
		return
	}
	if reply.Success == true {
		svr.ClientNumNotConMstr ++
	}
	
}

func (svr *Server)Check() bool{
	svr.ClientNumNotConMstr = 0
	for i := 0; i < len(svr.ClientList); i++ {
		fmt.Println(svr.ClientList[i].LocalPort)
		svr.CheckIfMasterDown(svr.ClientList[i])
	}
	if svr.ClientNumNotConMstr == len(svr.ClientList){
		fmt.Println(svr.ClientNumNotConMstr)
		fmt.Println(len(svr.ClientList))
		return true
	}
	return false
}

// RPC related
func (svr *Server)Operate(args *C2SArgs, reply *C2SReply) error{
	if !svr.isMaster{
		reply.Msg=WRONGLEAD
		reply.Success=false
		return nil
	}
	svr.update(args, reply)

	return nil
}

func (svr *Server) Exit(args C2SArgs, reply *C2SReply)error{
	//read seen table
	svr.mu.Lock()
	defer svr.mu.Unlock()
	fmt.Println(args.ClientId)
	//delete client from onlineclientlist
	delete(svr.OnlineClientList,args.ClientId)
	fmt.Printf("OnlineCltList length: %v\n",len(svr.OnlineClientList))
	//update seen table
	keysToUpdate := []string{}
	for filenameclientid, _ := range svr.SeenClientRid {
		strs:= strings.Split(filenameclientid,"//")
		if strs[1] == strconv.Itoa(args.ClientId){
			svr.SeenClientRid[filenameclientid] = SeenRidRecord{FilenameClient: filenameclientid, Rid: 0}
			keysToUpdate = append(keysToUpdate, filenameclientid)
		}
	}

	fmt.Printf("Rid with this client id: %v\n",svr.SeenClientRid["test.txt12//23"])
	fmt.Println(keysToUpdate)
	//send update to backup
	
	m2bArgs :=M2BArgs{UpdateType: UPDATETYPE5, OfflineClientId: args.ClientId, NewSeenRidRecordList:keysToUpdate}
	m2bReply := &M2BReply{}

	svr.sendRecordToBackup(svr.peer.IP, svr.peer.Port, m2bArgs, m2bReply)
	reply.Success=true

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
		//fmt.Println("dialing error")
		return false
	}

	err = client.Call("Client.Operate", args,reply)

	if err != nil{
		//fmt.Println("RPC failed")
		fmt.Println(err)
		return false
	}

	return true	
}


func (svr *Server) printTables(){
	fmt.Println("-------------------我是分割线---------------------")
	//FileOp
	fmt.Println("----------------FileOpTable")
	counter := 0
	for filename, record := range svr.FileOpTable{
		fmt.Printf("Record # %v\n", counter)
		fmt.Printf("filename: %v \n", filename)
		fmt.Printf("TimeStamp: %v\n", record.TimeStamp)
		fmt.Println("Status: "+record.Status)
		fmt.Printf("ClientList: %v \n",record.ClientIdList)
		counter += 1
	}


	//ClientO
	fmt.Println("----------------ClientOpTable")
	counter = 0
	for fc, record := range svr.ClientOpTable{
		fmt.Printf("Record # %v\n", counter)
		fmt.Printf("Clientid: %v\n", fc)
		fmt.Printf("FileList: %v \n",record)
		counter += 1
	}

	//SeenRid
	fmt.Println("-----------------SeenRid")
	counter = 0
	for _,record :=range svr.SeenClientRid{
		fmt.Printf("Record # %v\n", counter)
		fmt.Println("FilenameClient: "+record.FilenameClient)
		fmt.Printf("Rid: %v\n", record.Rid)
		counter += 1
	}

}