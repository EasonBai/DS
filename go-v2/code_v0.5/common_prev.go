

const (
	Create = 1
	Update = 2
	Delete = 3
)

type CtoSArgs struct{
	FileId		string
	OpType		int
	Content 	[]byte
}

type StoCArgs struct{
	FileId		string
	OpType		int 
	TimeStamp 	int64
	Content 	[]byte
}



