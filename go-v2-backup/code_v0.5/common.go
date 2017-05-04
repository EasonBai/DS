package hahaha

const (
	UPDATE = "UPDATE"
	DELETE = "DELETE"
)

type ClientEnd struct{
	//待修改
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
//FileOpTable	 map[string]FileOpRecord
type FileOpRecord struct{
	Filename 	string
	TimeStamp	int
	Status		string
	ClientList	[]int
	Content		[]byte
}
//ClientOpTable map[int]ClientOpRecord
type ClientOpRecord struct{
	ClientId	int
	FileList	[]string
}

type C2SArgs struct{
	ClientId 	int
	ReqId 		int
	Filename	string
	OpType		string
	Content 	[]byte
}

type C2SReply struct{
	Success string
}

type S2CArgs struct{
	ServerId	int
	FileId		string
	OpType		string 
	TimeStamp 	int
	Content 	[]byte
}

type S2CReply struct{
	Success string
}