package main

import (
	"fmt"
	"io/ioutil"

)
func main(){
	success, dat :=readFile("/Users/LIHUI/Documents/project651/test.txt")
	if success{
		writeFile("/Users/LIHUI/Documents/project651/hahaha.txt", dat)
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

