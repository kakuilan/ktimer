package main

import (
	. "ktimer"
    "fmt"
)

func main() {
	//命令行处理
    fmt.Println(111, RunLoger)
    ll,err := GetRunLoger()
    fmt.Println(222, ll,err)
    CatchCli()
}
