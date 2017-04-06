package ktimer

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"murmur3"
	"regexp"
	"strconv"
	"strings"
	"time"
	"github.com/go-redis/redis"
	"github.com/shopspring/decimal"
)

const (
	SERNAME = "ktimer"
	SERDESC = "Ktimer Service"
	PRODESC = "Ktimer is a simple timer/ticker task manager by golang."
	VERSION = "0.0.1"
	PUBDATE = "2017.4"
	AUTHOR  = "kakuilan@163.com"
    LOCKTIME = 2 * time.Second
)

//定时器参数数据结构
type KtimerData struct {
	Type    string `json:"type"`
	Time    int    `json:"time"`
	Limit   int    `json:"limit"`
	Command string `json:"command"`
}

//定时任务详情结构
type KtimerTask struct {
	KtimerData
	Run_num      int     `json:"run_num"`
	Run_lasttime float64 `json:"run_lasttime"`
	Run_nexttime float64 `json:"run_nexttime"`
}

//定时器容器
func TimerContainer() {
	go func() {
        //500毫秒的断续器
		mt := time.Tick(time.Millisecond * 500)
		for c := range mt {
			pidno, _ := GetServicePidNo()
			servpidno := GetCurrentServicePid()
			if pidno > 0 && pidno != servpidno {
				msg := fmt.Sprintf("check pid exception,service [%d] stopped.", servpidno)
				LogErres(msg)
			}

			//MainTimer()
            fmt.Println(c)
            _,now_mic := GetCurrentTime()
            LogRunes("ticker-MainTimer", now_mic)
            go func(now_mic float64) {
                runNum,runErr := MainTimer(now_mic)
                if runErr!=nil {
                    LogErres("MainTimer run error,", runErr)
                }else{
                    LogRunes("MainTimer run [%g] tasks at [%d]", runNum, now_mic)
                }
            }(now_mic)

		}
	}()
}

//主体定时器
func MainTimer(now_mic float64) (int,error) {
	var num int
	var err error
    var breakQue bool

    ms := GetMainSecond(now_mic)
    cnfObj, _ := GetConfObj()
    prefix := cnfObj.String("task_trun_key")
    key := prefix + strconv.Itoa(ms)

    client,err := GetRedisClient()
    if err!=nil {
        LogErres("MainTimer redis error", err)
    }

    taskExpire,err := cnfObj.Int("task_expire_limit")
    if err!=nil {
        LogErres("MainTimer conf task_expire_limit err", err)
    }else{
        taskExpire = 60
    }

    i := 0
    for {
        if breakQue {
            break
        }
        
        i++
        if i>5 {
            breakQue = true
        }

        zres,err := client.ZRangeWithScores(key, 0, 0).Result()
        zlen := len(zres)
        if err!=nil || zlen ==0 {
            breakQue = true
        }else{
            item := zres[0]
            if Greater(item.Score, now_mic) { //未到执行时间
                breakQue = true
            }else if Greater((now_mic- float64(taskExpire)), item.Score) { //过期,丢弃不执行
               // _ = client.ZRem(key, item.Member).Err()
               LogRunes("task is expired,deleted.", item)
            }else{ //执行任务
                runRes,runErr := RunSecondTask(item, now_mic)
                if runErr==nil && runRes {
                    num++
                }
            }

                runRes,runErr := RunSecondTask(item, now_mic)
                if runErr==nil && runRes {
                    num++
                }
            fmt.Printf("type: %T\n", item)
        }
        fmt.Printf("%+v", zres, taskExpire)
    }

    fmt.Println(breakQue,now_mic, ms)
    msg := fmt.Sprintf("MainTimer:%0.6f, run tasks total:%d", now_mic, num)
    LogRunes(msg)

	return num, err
}

//执行定时器秒任务
func RunSecondTask(zd redis.Z, now_mic float64) (bool,error) {
    var res bool
    var err error
    
    msg := fmt.Sprintf("begin single SecondTask:%0.6f kid:[%s] time:[%0.6f]", now_mic, zd.Member, zd.Score)
    LogRunes(msg)



    return res,err
}

//加入定时器
func AddTimer(td *KtimerData) (bool, string, *KtimerTask, error) {
	var res bool
	var err error
    var kid string
	var kt = &KtimerTask{}

	if td.Command == "" {
		err = errors.New("command is empty")
		return res, kid, kt, err
	}
 
	if td.Type != "timer" && td.Type != "ticker" {
		err = errors.New("type is error")
		return res, kid, kt, err
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
	kt.Type = td.Type
	kt.Time = td.Time
	kt.Limit = td.Limit
	kt.Command = td.Command

	now_sec, now_mic := GetCurrentTime()
	maxSeconds, maxTimestamp, err := GetSysTimestampLimit()
	if err != nil {
		err = errors.New("conf task_max_day is error")
		return res, kid, kt, err
	} else if td.Time <= maxSeconds {
		kt.Run_nexttime = float64(td.Time) + now_mic
	} else if td.Time > maxSeconds && td.Time < now_sec {
		err = errors.New("time as second cannot >" + strconv.Itoa(maxSeconds))
		return res, kid, kt, err
	} else if td.Time >= now_sec && td.Time <= maxTimestamp {
		kt.Run_nexttime = float64(td.Time)
	} else {
		err = errors.New("time as timestamp cannot>" + strconv.Itoa(maxTimestamp))
		return res, kid, kt, err
	}

	_, kid = MakeTaskKey(td.Command, td.Type, td.Time)
	secNum := GetMainSecond(kt.Run_nexttime)
	jsonRes, err := json.Marshal(*kt)

	res, err = _addTask2Pool(kid, jsonRes)
	if err != nil {
		_, _ = _delTask4Pool(kid)
		return res, kid, kt, err
	}
	res, err = _addTask2Queu(kid, kt.Run_nexttime, secNum)
	if res {
		LogRunes("add new task", kt)
	}

	return res, kid, kt, err
}

//删除定时器
func DelTimer(kid string) (bool, error) {
	res, err := DelTaskDetail(kid)
	if res {
		LogRunes("del a task", kid)
	}

	return res, err
}

//添加任务到任务池
func _addTask2Pool(kid string, task []byte) (bool, error) {
	var res bool
	var err error

	cnfObj, _ := GetConfObj()
	key := cnfObj.String("task_pool_key")
	client, err := GetRedisClient()
	if err != nil {
		return res, err
	}

	err = client.HSet(key, kid, string(task)).Err()
	if err == nil {
		res = true
	}

	return res, err
}

//从任务池删除任务
func _delTask4Pool(kid string) (bool, error) {
	var res bool
	var err error

	cnfObj, _ := GetConfObj()
	key := cnfObj.String("task_pool_key")
	client, err := GetRedisClient()
	if err != nil {
		return res, err
	}

	err = client.HDel(key, kid).Err()
	if err == nil {
		res = true
	}

	return res, err
}

//新增任务到运转队列
func _addTask2Queu(kid string, nextime float64, secondkey int) (bool, error) {
	var res bool
	var err error

	cnfObj, _ := GetConfObj()
	prefix := cnfObj.String("task_trun_key")
	key := prefix + strconv.Itoa(secondkey)
	client, err := GetRedisClient()
	if err != nil {
		return res, err
	}

	zd := redis.Z{nextime, kid}
	err = client.ZAdd(key, zd).Err()
	if err == nil {
		res = true
	}

	return res, err
}

//从运转队列删除任务
func _delTask4Queu(kid string, nextime float64) (bool, error) {
	var res bool
	var err error

	cnfObj, _ := GetConfObj()
	prefix := cnfObj.String("task_trun_key")

	secondkey := GetMainSecond(nextime)
	key := prefix + strconv.Itoa(secondkey)
	client, err := GetRedisClient()
	if err != nil {
		return res, err
	}

	err = client.ZRem(key, kid).Err()
	if err == nil {
		res = true
	}

	return res, err
}

//重新添加定时器任务
func ReaddTimer() {

}

//更新定时器
func UpdateTimer(oldkid string, kd *KtimerData) (bool, string, *KtimerTask, error) {
    var res bool
	var newkid string
	var kt = &KtimerTask{}
	var err error

	oldkid = strings.TrimSpace(oldkid)
	if oldkid == "" {
		err = errors.New("kid is empty")
		return res, newkid, kt, err
	} else if !IsNumeric(oldkid) {
		err = errors.New("kid is not numeric")
		return res, newkid, kt, err
	}

    //检查新的数据
    if kd.Type=="" && kd.Time==0 && kd.Limit==0 && kd.Command=="" {
        err = errors.New("no parameters to update")
        return res, newkid, kt, err
    }

    kt,err = GetTaskDetail(oldkid)
    if err!=nil {
        return res, newkid, kt, err
    }

    if kd.Type =="" {
        kd.Type = kt.Type
    }
    if kd.Time == 0 {
        kd.Time = kt.Time
    }
    if kd.Limit == 0 {
        kd.Limit = kt.Limit
    }
    if kd.Command == "" {
        kd.Command = kt.Command
    }

    //是否会生成新的任务
    _,newkid = MakeTaskKey(kd.Command, kd.Type, kd.Time)
    if oldkid!=newkid {
        _,_ = DelTimer(oldkid)
    }else{
        _,_ = _delTask4Queu(oldkid, kt.Run_nexttime)
    }

    res,newkid,kt,err = AddTimer(kd)
    if err!=nil {
        return res, newkid, kt, err
    }
 
	return res, newkid, kt, err
}

//获取定时器
func GetTimer(kid string) (*KtimerTask, error) {
	res, err := GetTaskDetail(kid)
	return res, err
}

//统计定时器任务
func CountTimer() (int, error) {
	var res int
	var err error

	cnfObj, _ := GetConfObj()
	key := cnfObj.String("task_pool_key")
	client, err := GetRedisClient()
	if err != nil {
		return res, err
	}

	num, err := client.HLen(key).Result()
	res = int(num)
	return res, err
}

//清空所有定时器
func ClearTimer() (bool, error) {
	var res bool
	var err error

	cnfObj, _ := GetConfObj()
	pool_key := cnfObj.String("task_pool_key")
	trun_prefix := cnfObj.String("task_trun_key")
	client, err := GetRedisClient()
	if err != nil {
		return res, err
	}

	err = client.Del(pool_key).Err()
	if err == nil {
		res = true
		for i := 0; i <= 59; i++ {
			trun_key := trun_prefix + strconv.Itoa(i)
			_ = client.Del(trun_key)
		}
	}

	return res, err
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

//获取主秒数
func GetMainSecond(t interface{}) int {
	var res int
	var mst int64
	var dec decimal.Decimal

	switch t.(type) {
	case int:
		tmp, _ := t.(int)
		mst = int64(tmp)
		dec = decimal.New(mst, 0)
	case int32:
		tmp, _ := t.(int32)
		mst = int64(tmp)
		dec = decimal.New(mst, 0)
	case int64:
		tmp, _ := t.(int64)
		mst = int64(tmp)
		dec = decimal.New(mst, 0)
	case float32:
		tmp, _ := t.(float32)
		flo := float64(tmp)
		dec = decimal.NewFromFloat(flo)
	case float64:
		flo, _ := t.(float64)
		dec = decimal.NewFromFloat(flo)
	case string:
		tmp, _ := t.(string)
		dec, _ = decimal.NewFromString(tmp)
	default:
		res = 1
	}

	if res != 1 {
		rem := dec.Mod(decimal.New(60, 0)).IntPart()
		res = int(rem)
	}

	return res
}

//获取当前时间
func GetCurrentTime() (int, float64) {
	var sec int
	var mic float64

	sec = int(time.Now().Unix())
	pn := math.Pow10(9)
	dic_pn := decimal.NewFromFloat(pn)
	dic_ms := decimal.New(time.Now().UnixNano(), 0)
	dic_ms = dic_ms.DivRound(dic_pn, 6)
	mic, _ = dic_ms.Float64()

	return sec, mic
}

//生成定时器ID
func MakeTimerId(command string, ttype string, ttime int) uint32 {
	var id uint32
	str := ttype + command
	if ttype == "timer" {
		str = str + strconv.Itoa(ttime)
	}

	key := []byte(str)
	id = murmur3.Sum32(key)

	return id
}

//生成定时任务key
func MakeTaskKey(command string, ttype string, ttime int) (uint32, string) {
	var id uint32
	var key string

	id = MakeTimerId(command, ttype, ttime)
	key = strconv.Itoa(int(id))

	return id, key
}

//检查字符串是否URL
func IsUrl(str string) bool {
	var res bool
	reg, _ := regexp.Compile(`^[a-zA-Z]+://(\w+(-\w+)*)(\.(\w+(-\w+)*))*(\?\s*)?$`)
	res = reg.Match([]byte(str))
	return res
}

//检查字符串是否数值
func IsNumeric(str string) bool {
	var res bool
	reg, _ := regexp.Compile(`^[0-9]+(.[0-9]*)?$`)
	res = reg.Match([]byte(str))
	return res
}

//获取任务详情
func GetTaskDetail(kid string) (*KtimerTask, error) {
	var err error
	var kd = &KtimerTask{}

	kid = strings.TrimSpace(kid)
	if kid == "" {
		err = errors.New("kid is empty")
		return kd, err
	} else if !IsNumeric(kid) {
		err = errors.New("kid is not numeric")
		return kd, err
	}

	cnfObj, _ := GetConfObj()
	key := cnfObj.String("task_pool_key")
	client, err := GetRedisClient()
	if err != nil {
		return kd, err
	}

	res, err := client.HGet(key, kid).Result()
	if err == nil {
		json.Unmarshal([]byte(res), kd)
	} else {
		err = errors.New("kid does not exist")
	}
	//fmt.Println(res,err)

	return kd, err
}

//删除任务详情
func DelTaskDetail(kid string) (bool, error) {
	var err error
	var res bool
	var kd = &KtimerTask{}

	kid = strings.TrimSpace(kid)
	if kid == "" {
		err = errors.New("kid is empty")
		return res, err
	} else if !IsNumeric(kid) {
		err = errors.New("kid is not numeric")
		return res, err
	}

	cnfObj, _ := GetConfObj()
	key := cnfObj.String("task_pool_key")
	client, err := GetRedisClient()
	if err != nil {
		return res, err
	}

	str, err := client.HGet(key, kid).Result()
	if err != nil {
		err = errors.New("kid does not exist")
		return res, err
	}

	err = json.Unmarshal([]byte(str), kd)
	if err != nil {
		return res, err
	}

	res, err = _delTask4Pool(kid)
	_, _ = _delTask4Queu(kid, kd.Run_nexttime)

	return res, err
}

//获取任务处理锁
func GetTaskDoingLock(kid string) (bool,error) {
    var res bool
    var err error

    if kid=="" {
        err = errors.New("kid is empty")
        return res,err
    }

    cnfObj, _ := GetConfObj()
    prefix := cnfObj.String("task_lcok_key")
    key := prefix + kid

    client, err := GetRedisClient()
    if err!=nil {
        return res,err
    }

    err = client.SetNX(key, 1, LOCKTIME).Err()
    if err==nil {
        res = true
    }

    return res,err
}

//解锁任务处理
func UnlockTaskDoing(kid string) (bool,error) {
    var res bool
    var err error

    if kid=="" {
        err = errors.New("kid is empty")
        return res,err
    }

    cnfObj, _ := GetConfObj()
    prefix := cnfObj.String("task_lcok_key")
    key := prefix + kid

    client, err := GetRedisClient()
    if err!=nil {
        return res,err
    }

    err = client.Del(key).Err()
    if err==nil {
        res = true
    }

    return res,err
}

//浮点数比较大小
func Greater(a,b float64) bool {
    return math.Max(a,b) ==a && math.Abs(a-b) > 0.000001
}

