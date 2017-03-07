package main

import (
	"fmt"
    "log"
	. "ktimer"
)

func main() {
	//命令行处理
	fmt.Println(111)
	rl, _ := GetRunLoger()
	el, _ := GetErrLoger()
	i := 20
    gg := log.New(rl,"",log.LstdFlags)
    for i > 0 {
		rl.Write([]byte("this is a runing log"))
		el.Write([]byte("this is a error log"))
	    gg.Println("system logger hahah")
        i--
	}
	fmt.Println(rl, el)
	CatchCli()
}
