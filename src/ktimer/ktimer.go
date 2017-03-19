package ktimer

import (
    "fmt"
    "time"
)

const (
    SERNAME = "ktimer"
    SERDESC = "Ktimer Service"
    PRODESC = "Ktimer is a simple timer/ticker task manager by golang."
	VERSION = "0.0.1"
    PUBDATE = "2017.3"
	AUTHOR  = "kakuilan@163.com"
)

//定时器数据结构
type KtimerData struct {
    Type string `json:"type"`
    Time int `json:"time"`
    Limit int `json:"limit"`
    Command string `json:"command"`
}

//定时任务详情结构
type KtaskDetail struct {
    KtimerData
    Run_num int `json:"run_num"`
    Run_lasttime  float32 `json:"run_lasttime"`
    Run_nexttime float32 `json:"run_nexttime"`
}


//定时器容器
func TimerContainer() {
    rlg,_ := GetRunLoger()
    elg,_ := GetErrLoger()
    go func(){
        
    mt := time.Tick(time.Millisecond * 500)
    for c := range mt {
        pidno,_ := GetServicePidNo()
        servpidno := GetCurrentServicePid()
        if pidno!=servpidno{
            msg := fmt.Sprintf("check pid exception,service [%d] stopped.", servpidno)
            elg.Println(msg)
        }

        now := time.Now().UnixNano()
        fmt.Println(mt,c, now)
        rlg.Println(mt, c, now)
        rlg.Println(SERNAME, "定时器运行")
        MainTimer()
    }
    }()
    
}

//主体定时器
func MainTimer() (bool,error) {
    var res bool = false
    var err error

    return res,err
}

//加入定时器
func AddTimer() {

}

//更新定时器
func UpdateTimer() {
    
}

//获取定时器
func GetTimer() {

}

//删除定时器
func DelTimer() {

}

//统计定时器
func CountTimer() {

}

//清空所有定时器
func ClearTimer() {

}

//执行定时器秒任务
func RunSecondTask() {
    
}

//执行具体任务
func RunDetailTask() {
    
}


