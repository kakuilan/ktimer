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
 
    //无参数直接显示使用方式
    if(argNum==1) {
        Help()
    }else{
        CatchCli()
    }


    //配置文件
    //file := GetConfFilePath()
    //fmt.Println(file)

    //初始化
    Init()


}
