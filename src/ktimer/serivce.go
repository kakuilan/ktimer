package ktimer

import (
	"fmt"
    //"github.com/astaxie/beego/config"
    //"config"
    "gopkg.in/redis.v5"
    "errors"
    "os"
    "strings"
)

//获取redis连接
func GetRedisClient() (*redis.Client,error) {
    var client *redis.Client
    var err error
    CnfObj,err = GetConfObj()

    if(err ==nil ) {
      var addr string
      host := CnfObj.String("redis::redis.host")
      port := CnfObj.String("redis::redis.port")
      pawd := CnfObj.String("redis::redis.passwd")
      db,err2 := CnfObj.Int("redis::redis.db")
      addr = host + ":" + port
      fmt.Println(host,port,addr,pawd,db,err2, "redis conf")
      if(err2!=nil) {
          return client,err2
      }

      client = redis.NewClient(&redis.Options{
          Addr: addr ,
          Password : pawd,
          DB: db,
      })

      return  client,err2
    }

    return client,err
}

//检查redis是否连接
func CheckRedis() (bool,error){
    var client *redis.Client
    var err error
    var pong string
    var res bool = false
    
    client,err = GetRedisClient()
    if(err !=nil ) {
        return  res,err
    }

    pong,err = client.Ping().Result()
    if(err !=nil ) {
        return res,err
    }else if(pong!="PONG") {
        err = errors.New("reids ping result not eq `PONG`")
        return res,err
    }

    return true,err
}

func CheckLogdir() (bool,error) {
    var err error
    CnfObj,err = GetConfObj()
    if(err !=nil ){
        return  false,err
    }

    logdir := CnfObj.String("log::log.dir")
    pos := strings.Index(logdir, "/")
    if(pos==-1) { //相对当前目录
        currdir := GetCurrentDirectory()
        logdir = currdir + "/" + logdir
    }

    write := Writeable(logdir)

    fmt.Println(logdir, pos, write)

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
    redisChk,err := CheckRedis()
    if(err!=nil){
        fmt.Println("redis connet has error:",  redisChk, err) 
        os.Exit(0)
    }

    //检查日志目录
    logChk,err2 := CheckLogdir()
    fmt.Println(logChk,err2)

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

