package main
import (
    . "ktimer"
    "fmt"
)

func main() {
    //配置文件
    cnf := GetConfFilePath()
    fmt.Println(cnf)

    //初始化
    Init()

}
