package ktimer

import (
    "fmt"
    "time"
)

const (
    PRODESC = "Ktimer is a simple timer/ticker manager by golang."
	VERSION = "0.0.1"
    PUBDATE = "2017.3"
	AUTHOR  = "kakuilan@163.com"
)

//定时器容器
func TimerContainer() {
    mt := time.Tick(time.Millisecond * 500)
    for c := range mt {
        now := time.Now().UnixNano()
        fmt.Println(mt,c, now)
        //MainTimer(c)
    }
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


