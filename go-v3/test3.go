package main
import (
	"bufio"
	"fmt"
	"os"
	//"io/ioutil"
	"strings"
	//"strconv"
)


func main() {


	//servers := make([]*Server, 2)

	// clientslist :=make(map[int]ClientEnd)
	// clientslist[0] = ClientEnd{Id: 0, IP: "127.0.0.1", LocalPort: "8000"}
	//clientslist[1] = ClientEnd{Id: 1, IP: "127.0.0.1", LocalPort: "8001"}

	serverslist :=make(map[int]ServerEnd)
	serverslist[0] = ServerEnd{Id: 0, IP:"127.0.0.1", LocalPort:"7998"}
	serverslist[1] = ServerEnd{Id: 1, IP:"127.0.0.1", LocalPort:"7999"}

	//paFor0 := PeerAddress{IP:"127.0.0.1", Port:"7999"}
	//paFor1 := PeerAddress{IP:"127.0.0.1", Port:"7998"}
	//master := new(Server)
	//master.StartServer(clientslist,0,true,"7998",paFor0)
	//servers[1]=StartServer(clientslist,1,false,"7999",paFor1)
	
	client0 := new(Client)

	client0.MakeClient(serverslist,0,0,"8000")
}



func test() {
	//read
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please enter 'get <filename>' or 'send <filename>' to transfer the files \n")
	inputFromUser, _ := reader.ReadString('\n')
	inputFromUser=strings.TrimSpace(inputFromUser)
	Str:=strings.Split(inputFromUser, "\\")
	filename:=Str[len(Str)-1]
	fmt.Println(filename)
	// arrayOfCommands := strings.Split(inputFromUser, "\\")
	// for r:=range arrayOfCommands{
	// 	fmt.Println(arrayOfCommands[r])
	// }
	fmt.Println("hahaha")
	// St,_ := reader.ReadString('\n')
	// // St=strings.TrimSpace(St)
	// // fmt.Println(len(St))
	// St=St[:1]
	// op,err:=strconv.Atoi(St)
	// if err==nil{
	// 	fmt.Print(op)
	// } else {
	// 	fmt.Print(err)
	// }
	//var AllClientsId []int
	map1:=make(map[int]string)
	//map1[1]=FileOpRecord{}
	map1[1]="123"
	map1[2]="123"
	map1[3]="123"
	map1[4]="123"
	// fmt.Println(map1[5].TimeStamp)
	// fmt.Println(map1[1].TimeStamp)
	fmt.Println("end")
	map2:= make(map[int](map[string]struct{}))
	for key := range map1{
		map2[key]=make(map[string]struct{})
	}
	map2[1]["00"]=struct{}{}
	map2[1]["11"]=struct{}{}
	delete(map2[1],"00")
	delete(map2[1],"11")
	delete(map2[1],"11")
	if _,contains :=map2[1]; contains{
		// for id:=range map2[1]{
		// 	fmt.Println(id)
		// }
		fmt.Println(len(map2[1]))
	} else {
		fmt.Println("no")
	}
	// fpath:="G:\\BU\\distributed system\\project\\1.txt"
	// dat, err := ioutil.ReadFile(fpath)
	// fpath="temp"
	//err=os.Mkdir("src", 0777)
	// fpath="G:\\BU\\distributed system\\project\\code\\temp\\3.txt"
	//err = ioutil.WriteFile("src\\"+"3.txt", dat, 0644)
	// err:=os.Remove("src\\3.txt")
	// if err==nil{
	// 	fmt.Println("yeah!")
	// } else {
	// 	fmt.Println(err)
	// }
	// for  key := range map1 {
	// 	AllClientsId=append(AllClientsId, key)
	// }
	// fmt.Println(AllClientsId)
	//filename="args.Filename"
	// cid:=3
	// Rid:=2
	// key:=filename+"//"+strconv.Itoa(cid)
	// Entry := M2BArgs{}
	// Entry.UpdateType=UPDATETYPE1
	// Entry.NewSeenRidRecord=SeenRidRecord{FilenameClient: key, Rid: Rid}
	// //Entry.NewClientOpRecord=
	// fmt.Println(Entry.UpdateType)
	// fmt.Println(Entry.NewFileOpRecord)
	//inputFromUser, _ = reader.ReadString('\n')
	// inputFromUser = "huihiu"
	// fmt.Println(inputFromUser+"ttt")
}