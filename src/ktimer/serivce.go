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

//检查日志目录
func CheckLogdir() (bool,error) {
    var err error
    CnfObj,err = GetConfObj()
    if(err !=nil ){
        return  false,err
    }

    logdir := CnfObj.String("log::log.dir")
    logdir = strings.Replace(logdir, "\\", "/", -1)
    pos := strings.Index(logdir, "/")
    if(pos==-1) { //相对当前目录
        currdir := GetCurrentDirectory()
        logdir = currdir + "/" + strings.TrimRight(logdir, "/")
    }

    direxis := FileExist(logdir)
    if(!direxis) {
        err = os.MkdirAll(logdir, 0766)
        if (err !=nil) {
            return false ,err
        }
    }else{
        write := Writeable(logdir)
        err = os.Chmod(logdir, 0766)
        if(!write || err!=nil) {
            err = errors.New("logdir canot write")  
            return false, err
        }
    }

    return true,err
}

//检查pid文件
func CheckPid() (bool,error) {
    var err error
    var chk bool = false

    return chk,err
}

//初始化
func Init() {
    var err error
    var chk bool
   
    fmt.Println("CnfObj", CnfObj)

    //检查配置文件
    chk = CheckConfFile()
    if(!chk) {
        CreateConfFile()
    }

    //检查redis
    chk,err = CheckRedis()
    if(err!=nil){
        fmt.Println("redis connet has error:", err) 
        os.Exit(0)
    }

    //检查日志目录
    chk,err = CheckLogdir()
    if(err !=nil) {
        fmt.Println("check log`s dir has error:", err)
    }

    //检查pid
    chk,err = CheckPid()


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

