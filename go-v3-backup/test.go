package main
import (
	//"bufio"
	"fmt"
	"os"
	//"io/ioutil"
	//"strings"
	//"bytes"
	//"strconv"
)
// func writeFile(path string, content []byte) bool{
// 	err := ioutil.WriteFile(path, content, 0644)

// 	if err != nil{
// 		fmt.Println(err)
// 		return false
// 	}

// 	return true
// }


func main() {
	//read
	// reader := bufio.NewReader(os.Stdin)
	// fmt.Print("Please enter 'get <filename>' or 'send <filename>' to transfer the files \n")
	// inputFromUser, _ := reader.ReadString('\n')
	// inputFromUser=strings.TrimSpace(inputFromUser)
	// Str:=strings.Split(inputFromUser, "\\")
	// filename:=Str[len(Str)-1]
	// fmt.Println(filename)
	// // arrayOfCommands := strings.Split(inputFromUser, "\\")
	// // for r:=range arrayOfCommands{
	// // 	fmt.Println(arrayOfCommands[r])
	// // }
	// fmt.Println("hahaha")
	// _,data:=readFile("G:\\BU\\distributed system\\project\\1.txt")
	// //err:=writeFile(filename, data)
	// _,data2:=readFile(filename)
	// file,err:=os.Open(filename)
	// err=file.Close()
	// fmt.Println(bytes.Compare(data,data2))
	// if err==nil{
	// 	fmt.Println("Success!")
	// } else {
	// 	fmt.Println("Failed...")
	// }
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
	// err:=os.Remove("src\\"+"1.txt")
	_, err := os.Stat("src\\"+"1.txt")
	fmt.Println(os.IsNotExist(err))
	//fmt.Println(err)
	map1:=make(map[int]string)
	//map1[1]=FileOpRecord{}
	map1[1]="123"
	map1[2]="123"
	map1[3]="123"
	map1[4]="123"
	map3:=make(map[int]string)
	map3=map1
	//delete(map3, 1)
	delete(map3,2)
	fmt.Println(len(map1))
	//fmt.Println(len(map3))
	// fmt.Println(map1[5].TimeStamp)
	// fmt.Println(map1[1].TimeStamp)
	fmt.Println("end")
	map2:= make(map[int](map[string]struct{}))
	for key := range map1{
		map2[key]=make(map[string]struct{})
	}
	map2[1]["00"]=struct{}{}
	if _,contains :=map2[1]["00"]; contains{
		fmt.Println("1")
	} else {
		fmt.Println("2")
	}
	map2[1]["11"]=struct{}{}
	delete(map2[1],"00")
	if _,contains :=map2[1]["00"]; !contains{
		fmt.Println("3")
	} else {
		fmt.Println("4")
	}
	map2[1]["00"]=struct{}{}
	if _,contains :=map2[1]["00"]; contains{
		fmt.Println("5")
	} else {
		fmt.Println("6")
	}
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
	if _,contains :=map2[1]["11"]; contains{
		fmt.Println("7")
	} else {
		fmt.Println("8")
	}
	if _,contains :=map2[1]["00"]; contains{
		fmt.Println("9")
	} else {
		fmt.Println("10")
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
	// filename="args.Filename"
	// cid:=3
	// Rid:=2
	// key:=filename+"//"+strconv.Itoa(cid)
	// Entry := M2BArgs{}
	// Entry.UpdateType=UPDATETYPE1
	// Entry.NewSeenRidRecord=SeenRidRecord{FilenameClient: key, Rid: Rid}
	// // //Entry.NewClientOpRecord=
	// fmt.Println(Entry.UpdateType)
	// fmt.Println(Entry.NewFileOpRecord)
	//inputFromUser, _ = reader.ReadString('\n')
	// inputFromUser = "huihiu"
	// fmt.Println(inputFromUser+"ttt")
}