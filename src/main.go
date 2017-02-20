package main
import (
    . "ktimer"
    "fmt"
)

func main() {
    //配置文件
    file := GetConfFilePath()
    fmt.Println(file)

    //初始化
    Init()
}
