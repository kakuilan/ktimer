package main
import (
    . "ktimer"
    "fmt"
    "os"
)

func main() {
    //获取命令行参数
    argNum := len(os.Args)
    fmt.Println(argNum)
    fmt.Println(os.Args)
   
    if(argNum==1) {
        Help()
    }


    //配置文件
    //file := GetConfFilePath()
    //fmt.Println(file)

    //初始化
    Init()


}
