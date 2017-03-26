package ktimer

import (
	"errors"
	"fmt"
	"murmur3"
	"time"
	"github.com/shopspring/decimal"
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
	Type    string `json:"type"`
	Time    int    `json:"time"`
	Limit   int    `json:"limit"`
	Command string `json:"command"`
}

//定时任务详情结构
type KtaskDetail struct {
	KtimerData
	Run_num      int     `json:"run_num"`
	Run_lasttime float32 `json:"run_lasttime"`
	Run_nexttime float32 `json:"run_nexttime"`
}

//定时器容器
func TimerContainer() {
	rlg, _ := GetRunLoger()
	elg, _ := GetErrLoger()
	go func() {

		mt := time.Tick(time.Millisecond * 500)
		for c := range mt {
			pidno, _ := GetServicePidNo()
			servpidno := GetCurrentServicePid()
			if pidno>0 && pidno != servpidno {
				msg := fmt.Sprintf("check pid exception,service [%d] stopped.", servpidno)
				elg.Println(msg)
			}

			now := time.Now().UnixNano()
			fmt.Println(mt, c, now)
			rlg.Println(mt, c, now)
			rlg.Println(SERNAME, "定时器运行")
			MainTimer()
		}
	}()

}

//主体定时器
func MainTimer() (bool, error) {
	var res bool = false
	var err error

	return res, err
}

//加入定时器
func AddTimer(td KtimerData) (bool, error) {
	var res bool
	var err error

	if td.Command == "" {
		err = errors.New("command is empty")
		return res, err
	}

	if td.Type != "timer" && td.Type != "ticker" {
		err = errors.New("type is error")
		return res, err
	}

	if td.Time <= 0 {
		td.Time = 5
	}

	if td.Limit < 0 && td.Type == "timer" {
		td.Limit = 1
	} else if td.Limit < 0 && td.Type == "ticker" {
		td.Limit = 0
	}

	//定时器详情
	detail := KtaskDetail{
		td,
		0,
		0.0,
		0.0,
	}
	detail.Time = 5

	maxSeconds, maxTimestamp, err := GetSysTimestampLimit()
	if err != nil {
		err = errors.New("conf task_max_day is error")
		return res, err
	}

	ts := time.Now().Unix()
	fmt.Println(detail, ts, maxSeconds, maxTimestamp)
	fmt.Printf("%T\n", ts)

	aa := 23412.443
	bb := GetMainSecond(aa)
	fmt.Println("res:", aa, bb)


	return res, err
}

func ReaddTimer() {

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

//从配置获取最大秒数
func GetSysTimestampLimit() (int, int, error) {
	var err error
	var maxSec, maxTim, maxDay int
	CnfObj, err = GetConfObj()
	if err != nil {
		return maxSec, maxTim, err
	}

	maxDay, err = CnfObj.Int("task_max_day")
	if err != nil {
		return maxSec, maxTim, err
	}

	maxSec = maxDay * 86400
	maxTim = int(time.Now().Unix()) + maxSec
	return maxSec, maxTim, err
}

//生成定时器ID
func MakeTimerId(command string) uint32 {
	key := []byte(command)
	id := murmur3.Sum32(key)
	return id
}

//获取主秒数
func GetMainSecond(t interface{}) int {
	var res int 
    var mst int64
	var dec decimal.Decimal

	switch t.(type) {
		case int :
			tmp,_ := t.(int)
			mst = int64(tmp)
			dec = decimal.New(mst, 1)
		case int32 :
			tmp,_ := t.(int32)
			mst = int64(tmp)
			dec = decimal.New(mst, 0)
		case int64 :
			tmp,_ := t.(int64)
			mst = int64(tmp)
			dec = decimal.New(mst, 0)
		case float32 :
			tmp,_ := t.(float32)
			flo := float64(tmp)
			dec = decimal.NewFromFloat(flo)
		case float64 :
			flo,_ := t.(float64)
			dec = decimal.NewFromFloat(flo)
		case string :
			tmp,_ := t.(string)
			dec,_ = decimal.NewFromString(tmp)
		default :
			res = 1
	}

	if res!=1 {
		rem := dec.Mod(decimal.New(60, 0)).IntPart()
		res = int(rem)
	}

	return res
}
