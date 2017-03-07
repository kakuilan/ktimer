package main

import (
	"fmt"
	. "ktimer"
	//"log"
)

func main() {
	//命令行处理
	fmt.Println(111)
	rl,e1 := GetRunLoger()
	el,e2 := GetErrLoger()
	i := 20
	for i > 0 {
		rl.Println("this is a runing log.")
		el.Println("this is a error log.")
		i--
	}
	fmt.Println(rl, el, e1, e2)
	CatchCli()
}
