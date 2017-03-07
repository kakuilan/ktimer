package ktimer
import (
    "gopkg.in/natefinch/lumberjack.v2"
    "log"
)

//全局日志对象
var RunLoger,ErrLoger lumberjack.Logger = nil

//获取运行日志对象
func GetRunLoger()(interface{}, error) {
    var err error
    var logdir string
    if RunLoger == nil {
        CnfObj, err = GetConfObj()
        if err !=nil {
            return  RunLoger,err
        }
        logdir,err = CheckLogdir()
        if err !=nil  {
           return RunLoger,err
        }
        
        file :=  CnfObj.String("log::log.runed_file")
        maxsize,err := CnfObj.Int("log::log.maxsize")
        if err!=nil {
            maxsize = 500
        }
        maxbackup,err := CnfObj.Int("log::log.maxbackup")
        if err !=nil {
            maxbackup = 5
        }
        maxage,err := CnfObj.Int("log::log.maxage")
        if err != nil {
            maxage = 30
        }

        RunLoger = &lumberjack.Logger{
            Filename : file,
            MaxSize : maxsize,
            MaxBackups : maxbackup,
            MaxAge : maxage,
        }
        println(RunLoger)
        //log.SetOutput(RunLoger) 
    }

    return RunLoger,err
}

//获取错误日志对象
func GetErrLoger(){
    
}
