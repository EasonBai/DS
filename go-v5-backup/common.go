package main

const (
	UPDATE = "UPDATE"
	DELETE = "DELETE"
	WRONGLEAD="WRONGLEADER"
	OLDREQ = "OLDREQUEST"
	NOTSUCCESS="NOTSUCCESS"
	ACCEPTED="ACCEPTED"
	BACKUPLOST="BACKUPLOST"
	//S2CReply
	DUPLICATE="DUPLICATEWRITE"
	WRITTEN="WRITTEN"
	UNABLEWRITE="UNABLETOWRITE"
	REMOVED="REMOVED"
	UNABLEREMOVE="UNABLETOREMOVE"
	//backupdate
	UPDATETYPE1 = 1 //Only Rid
	UPDATETYPE2 = 2 //Only Rid & FileOpRecord
	UPDATETYPE3 = 3 //Only FileOpRecord &ClientOpRecord; update 2 tables
	UPDATETYPE4 = 4 //Only FileOpRecord &ClientOpRecord delete whole record in backup
	UPDATETYPE5 = 5 //Client exits: delete client in onlineclientlist && update seenclientridtable
	//UPDATETYPE6 = 6 //Client enters: add client in onlineclientlist
)

type ClientEnd struct{
	Id 		int
	IP		string
	LocalPort 	string
	GroupID	string
}

type ServerEnd struct{
	Id int
	IP string
	LocalPort string
}

type PeerAddress struct{
	IP string
	Port string
}

//FileOpTable	 map[string]FileOpRecord
type FileOpRecord struct{
	Filename 	string
	TimeStamp	int
	Status		string
	ClientIdList	map[int]struct{}
	Content		[]byte
}

type ClientOpRecord struct{
	ClientId	int
	FileIdList	map[string]struct{}
}

type SeenRidRecord struct{
	FilenameClient string
	Rid int
}

type C2SArgs struct{
	ClientId 	int
	ReqId 		int
	Filename	string
	OpType		string
	Content 	[]byte
	ClientIp string
	ClientPort string
}

type C2SReply struct{
	Success bool
	Msg 	string
}

type S2CArgs struct{
	ClientId	int
	FileId		string
	OpType		string 
	TimeStamp 	int
	Content 	[]byte
	AllClients map[int]struct{}
	OnlineClients map[int]struct{}
	NewMasterId int
}

type S2CReply struct{
	Success bool
	Msg		string
}

type M2BArgs struct{
	UpdateType			int
	NewFileOpRecord		FileOpRecord
	NewClientOpRecord	ClientOpRecord
	NewSeenRidRecord 	SeenRidRecord
	NewSeenRidRecordList []string
	OfflineClientId int
	ClientId 	int
	ClientIp string
	ClientPort string
}

type M2BReply struct{
	Success bool
	TimeStamp int
}

type TestCaseArgs struct{
	Case int
	SleepDuration int
}

type TestCaseReply struct{
	Success bool
}