package main

import (
	"MyPHPServer/tool/file"
	"fmt"
	"io/ioutil"
)

func ListFile(ROOT, Path string) string {
	dir_list, e := ioutil.ReadDir(ROOT + Path)
	if e != nil {
		fmt.Println("read dir error")
		return "error"
	}
	/*
	   for i, v := range dir_list {
	       fmt.Println(i, "=", v.Name())
	   }*/
	tatol := "<hr>"
	for i, v := range dir_list {
		var myname string
		if file.IsDir(ROOT + Path + "/" + v.Name()) {
			myname = v.Name() + "/"
		} else {
			myname = v.Name()
		}

		tatol = tatol + "<a href='" + Path + "/" + v.Name() + "'>" + myname /*v.Name()*/ + "</a><hr>"
		//log.Print(i,string(rune(i)))
		fmt.Println(i, v.Name()+string(rune(i)))
	}
	//log.Print(tatol)
	return "<h1>Index Of</h1>" + tatol + "<p align='center'>Powered By: FlyKO</p>"

}