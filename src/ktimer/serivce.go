package ktimer

import (
	"fmt"
    "github.com/astaxie/beego/config"
)

//检查redis是否连接
func CheckRedis() (bool,error){
    file := GetConfFilePath()
    cnf,err := config.NewConfig("ini", file)
    if(err == nil) {
        host := cnf.String("redis.host")
        port := cnf.String("redis.post")
        db := cnf.String("redis.db")

        fmt.Println(host,port,db,"redis conf")
    }

    fmt.Println(cnf,err, "conf")

    return true,err
}


func CheckPid() {
}

//初始化
func Init() {
    //检查配置文件
    confCheck := CheckConfFile()
    if(!confCheck) {
        CreateConfFile()
    }

    //检查redis
    CheckRedis()

}

func Start() {
    
}

func Stop(){
    
}

func Restart(){
    
}

func Add(){
    
}

func Get(){
    
}

func Del(){
    
}

func Count(){
    
}

func List(){
    
}

func Clear() {
    
}

func Info() {
    
}

func Status() {
    
}

