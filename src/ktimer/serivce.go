package ktimer

import (
	"fmt"
    //"github.com/astaxie/beego/config"
    "config"
    "gopkg.in/redis.v5"
)

//
func GetRedisClient() (*redis.Client,error) {
    var client *redis.Client
    file := GetConfFilePath()
    cnf,err := config.NewConfig(file)
    if(err ==nil ) {
      var addr string
      host := cnf.String("redis::redis.host")
      port := cnf.String("redis::redis.port")
      pawd := cnf.String("redis::redis.passwd")
      db,err2 := cnf.Int("redis::redis.db")
      addr = host + ":" + port
      fmt.Println(host,port,addr,pawd,db,err2, "redis conf")
      
      client = redis.NewClient(&redis.Options{
          Addr: addr ,
          Password : pawd,
          DB: db,
      })

      fmt.Println(client)
      return  client,err2
    }else{
    }

    return client,err
}

//检查redis是否连接
func CheckRedis() (bool,error){
    client,err := GetRedisClient()
    fmt.Println(client)

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

