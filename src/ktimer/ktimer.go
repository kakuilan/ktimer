package ktimer

import (
	"encoding/json"
	"errors"
	"fmt"
	curl "github.com/andelf/go-curl"
	"github.com/go-redis/redis"
	"github.com/shopspring/decimal"
	"math"
	"murmur3"
	"net/url"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
//    "os"
)

const (
	SERNAME     = "ktimer"
	SERDESC     = "Ktimer Service"
	PRODESC     = "Ktimer is a simple timer/ticker task manager by golang."
	VERSION     = "0.0.1"
	PUBDATE     = "2017.4"
	AUTHOR      = "kakuilan@163.com"
	LOCKTIME    = 2 * time.Second
	TASKMAXTIME = 30
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

//处理结果结构
type TkProceResp struct {
	Kid string
	Res bool
	Err error
}

//禁止的命令
var DenyCmd = []string{
	"rm ",
	"kil ",
	"mv ",
	"wget ",
	"dd ",
}

//定时器容器
func TimerContainer() {
	go func() {
		//500毫秒的断续器
		mt := time.Tick(time.Millisecond * 1000)
		for _ = range mt {
			pidno, _ := GetServicePidNo()
			servpidno := GetCurrentServicePid()
			if pidno > 0 && pidno != servpidno {
				msg := fmt.Sprintf("check pid exception,service [%d] stopped.", servpidno)
				LogErres(msg)
			}

			_, now_mic := GetCurrentTime()
			go func(now_mic float64) {
				msg := fmt.Sprintf("MainTimer begining: now[%0.6f]", now_mic)
				LogRunes(msg)
				//fmt.Println(msg)
				_, runErr := MainTimer(now_mic)
				if runErr != nil {
					LogErres("MainTimer run error,", runErr)
				}
			}(now_mic)

		}
	}()
}

//主体定时器
func MainTimer(now_mic float64) (int, error) {
	var sucNum int
    var allNum int64
	var err error
	var breakQue bool
	var redZ redis.Z
	var redZArr []redis.Z

	ms := GetMainSecond(now_mic)
	cnfObj, _ := GetConfObj()
	prefix := cnfObj.String("task_trun_key")
	key := prefix + strconv.Itoa(ms)

	client, err := GetRedisClient()
	if err != nil {
		LogErres("MainTimer redis error", err)
		return sucNum, err
	}

	for {
		if breakQue {
			break
		}
		zres, err := client.ZRangeWithScores(key, allNum, allNum).Result()
		zlen := len(zres)
		if err != nil || zlen == 0 {
            //fmt.Println("queue is to end")
			break
		} else {
			if redZ == zres[0] {
                //fmt.Println("redZ=zres[0]", redZ, zres, allNum)
				break
			}

			allNum++
			redZ = zres[0]
			zms := GetMainSecond(redZ.Score)
            fmt.Println("redZ-data", redZ, ms, zms)
			if ms != zms && GreaterOrEqual(redZ.Score, now_mic) { //未到执行时间
				breakQue = true
				msg := fmt.Sprintf("not run time, nowtime[%0.6f] nextime[%0.6f] item:%v", now_mic, redZ.Score, redZ)
				LogRunes(msg)
				//fmt.Println(msg)
                break
			} else { //执行任务
                redZArr = append(redZArr, redZ)
                //fmt.Println("append", redZArr)
			}
		}
	}

	//channel
	ch := make(chan *TkProceResp, 100)
    chNum := len(redZArr)
    fmt.Println("total channel", allNum, chNum, redZArr)

    for _, redZ = range redZArr {
        //fmt.Println("add redZArr", redZ)
        //LogRunes("paifa channel", redZ)
		go func(zd redis.Z, now_mic float64, ch chan *TkProceResp) {
			runRes, runErr := RunSecondTask(zd, now_mic, ch)
			if !runRes || runErr != nil {
				_delTask4Queu(fmt.Sprintf("%s", redZ), redZ.Score)
			}
		}(redZ, now_mic, ch)
		//fmt.Printf("zstruct type:[%T] %v i[%d]\n", redZ, redZ, allNum)
	}

	//等待结果返回
	retNum := 0
	for {
		if retNum >= chNum {
			close(ch)
			break
		}

		select {
		case tpr := <-ch:
			retNum++
			if tpr.Res && tpr.Err == nil {
				sucNum++
			}
		default:
			time.Sleep(100 * time.Microsecond)
		}
	}

    if chNum>0 {
        //os.Exit(0)
    }

	msg := fmt.Sprintf("MainTimer result: time:%0.6f,second:%d tasks total:[%d] runed:[%d] chan:[%d] return:[%d]", now_mic, ms, allNum, sucNum, chNum, retNum)
	if allNum > 0 {
		LogRunes(msg)
	}
	//fmt.Println(msg)
	return sucNum, err
}

//执行定时器秒任务
func RunSecondTask(zd redis.Z, now_mic float64, ch chan *TkProceResp) (bool, error) {
	var res bool
	var err error
	var tpr = &TkProceResp{}

	kid := fmt.Sprintf("%v", zd.Member)
	tpr.Kid = kid
    LogRunes("SecondTask begining:", zd, now_mic)

	kd, err := GetTaskDetail(kid)
	if err != nil {
		delRes, delErr := DelTaskDetail(kid)
		if !delRes || delErr != nil {
			_, _ = _delTask4Queu(kid, zd.Score)
		}
		LogRunes("SecondTask is not exist,deleted.kid:", kid)
		//fmt.Println("kid not exist", delRes, delErr)

		tpr.Err = err
		ch <- tpr
		return res, err
	}

	//获取任务执行锁
	locked, _ := GetTaskDoingLock(kid)
	if !locked {
		LogRunes("SecondTask get doing lock faili.kid:", kid)
		err = errors.New("get doing lock fail")

		tpr.Err = err
		ch <- tpr
		return res, err
	}

	//达到执行次数限制,删除任务
	if kd.Type == "ticker" && kd.Limit > 0 && kd.Limit <= kd.Run_num {
		_, _ = UnlockTaskDoing(kid)
		_, _ = DelTaskDetail(kid)
		msg := fmt.Sprintf("SecondTask kid[%s] had runed [%d] times,deleted.", zd.Member, kd.Run_num)
		LogRunes(msg, kd)

		ch <- tpr
		return res, err
	}

	//检查是否过期
	cnfObj, _ := GetConfObj()
	taskExpire, err := cnfObj.Int("task_expire_limit")
	if err != nil {
		taskExpire, err = 60, nil
	}
	if kd.Type=="timer" && GreaterOrEqual(now_mic-float64(taskExpire), kd.Run_nexttime) {
		_, _ = UnlockTaskDoing(kid)
		_, _ = DelTaskDetail(kid)
		msg := fmt.Sprintf("SecondTask is expired. kid[%s] nowtime[%0.6f] nextime[%0.6f] expire[%d],deleted.", kid, now_mic, kd.Run_nexttime, taskExpire)
		LogRunes(msg)
		//fmt.Println(msg)
	}

	//执行日志
	msg := fmt.Sprintf("SecondTask begining:%0.6f kid[%s] time[%0.6f]", now_mic, zd.Member, zd.Score)
	LogRunes(msg)

	//删除任务
	_, _ = DelTaskDetail(kid)

	//若是ticker,重新加入
	if kd.Type == "ticker" && (kd.Limit == 0 || (kd.Limit > 0 && (kd.Limit-1) > kd.Run_num)) {
		_, addErr := ReaddTimerAfterRun(kid, kd)
		if addErr != nil {
			LogRunes("SecondTask readd task fail:", addErr)
		}
	}

	//执行
	RunDetailTask(kid, kd.Command)

	res, err = true, nil
	tpr.Res = true
	ch <- tpr
	return res, err
}

//执行具体任务
func RunDetailTask(kid string, command string) (bool, error) {
	var res bool
	var err error

	command = strings.TrimSpace(command)
	if IsUrl(command) { //执行URL任务
		out, err := RunUrlTask(command, true)
		if err == nil {
			res = true
		}
		LogRunes("exec url task res:", kid, command, out, err)
	} else { //命令行任务
		isDeny := false
		for _, cmd := range DenyCmd {
			if strings.Index(command, cmd) != -1 {
				isDeny = true
				err = errors.New("task contains deny cmd:" + cmd)
				break
			}
		}

		if !isDeny {
			out, err := RunCmdTask(command, false)
			if err == nil {
				res = true
			}
			LogRunes("exec cli task res:", kid, command, out, err)
		}
	}

	//解锁
	_, _ = UnlockTaskDoing(kid)

	return res, err
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

	if td.Type == "" {
		td.Type = "timer"
	} else if td.Type != "timer" && td.Type != "ticker" {
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
	if err != nil {
		err = errors.New("task detail json encode fail")
		return res, kid, kt, err
	}

	//检查任务数量限制
	cnfObj, _ := GetConfObj()
	tskMaxNum, err := cnfObj.Int("task_max_num")
	if err != nil || tskMaxNum <= 0 {
		tskMaxNum = 500000
	}
	curTskNum, err := CountTimer()
	if err != nil {
		return res, kid, kt, err
	} else if curTskNum >= tskMaxNum {
		err = errors.New("max number allowed task " + strconv.Itoa(tskMaxNum))
		return res, kid, kt, err
	}

	res, err = _addTask2Pool(kid, jsonRes)
	if err == nil {
		res, err = _addTask2Queu(kid, kt.Run_nexttime, secNum)
		if err != nil {
			_, _ = _delTask4Pool(kid)
		}
	}
	if res {
		LogRunes("add new task:", kid, kt)
	}

	return res, kid, kt, err
}

//执行后重新添加定时任务
func ReaddTimerAfterRun(kid string, kt *KtimerTask) (bool, error) {
	var res bool
	var err error

	kid = strings.TrimSpace(kid)
	if kid == "" {
		err = errors.New("kid is empty")
		return res, err
	} else if !IsNumeric(kid) {
		err = errors.New("kid is not numeric")
		return res, err
	}

	kt.Run_num++
	kt.Run_lasttime = kt.Run_nexttime
	if kt.Type != "ticker" || (kt.Limit > 0 && kt.Run_num >= kt.Limit) {
		err = errors.New("task is not ticker or number limit")
		return res, err
	}

	now_sec, now_mic := GetCurrentTime()
	maxSeconds, maxTimestamp, err := GetSysTimestampLimit()
	if err != nil {
		err = errors.New("conf task_max_day is error")
		return res, err
	} else if kt.Time <= maxSeconds {
		kt.Run_nexttime = float64(kt.Time) + now_mic
	} else if kt.Time > maxSeconds && kt.Time < now_sec {
		err = errors.New("time as second cannot >" + strconv.Itoa(maxSeconds))
		return res, err
	} else if kt.Time >= now_sec && kt.Time <= maxTimestamp {
		kt.Run_nexttime = float64(kt.Time)
	} else {
		err = errors.New("time as timestamp cannot>" + strconv.Itoa(maxTimestamp))
		return res, err
	}

	secNum := GetMainSecond(kt.Run_nexttime)
	jsonRes, err := json.Marshal(*kt)
	if err != nil {
		err = errors.New("task detail json encode fail")
		return res, err
	}

	res, err = _addTask2Pool(kid, jsonRes)
	if err == nil {
		res, err = _addTask2Queu(kid, kt.Run_nexttime, secNum)
		if err != nil {
			_, _ = _delTask4Pool(kid)
		}
	}

	if res {
		LogRunes("readd task:", kid, kt)
	}

	return res, err
}

//删除定时器
func DelTimer(kid string) (bool, error) {
	res, err := DelTaskDetail(kid)
	if res {
		LogRunes("del a task:", kid)
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
	if kd.Type == "" && kd.Time == 0 && kd.Limit == 0 && kd.Command == "" {
		err = errors.New("no parameters to update")
		return res, newkid, kt, err
	}

	kt, err = GetTaskDetail(oldkid)
	if err != nil {
		return res, newkid, kt, err
	}

	if kd.Type == "" {
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
	_, newkid = MakeTaskKey(kd.Command, kd.Type, kd.Time)
	if oldkid != newkid {
		_, _ = DelTimer(oldkid)
	} else {
		_, _ = _delTask4Queu(oldkid, kt.Run_nexttime)
	}

	res, newkid, kt, err = AddTimer(kd)
	if err != nil {
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
	reg, _ := regexp.Compile(`^http[s]?:\/\/(.*)(\?\s*)?`)
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
		LogErres("HGet kid err:", kid, err)
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
func GetTaskDoingLock(kid string) (bool, error) {
	var res bool
	var err error

	if kid == "" {
		err = errors.New("kid is empty")
		return res, err
	}

	cnfObj, _ := GetConfObj()
	prefix := cnfObj.String("task_lcok_key")
	key := prefix + kid

	client, err := GetRedisClient()
	if err != nil {
		return res, err
	}

	err = client.SetNX(key, 1, LOCKTIME).Err()
	if err == nil {
		res = true
	}

	return res, err
}

//解锁任务处理
func UnlockTaskDoing(kid string) (bool, error) {
	var res bool
	var err error

	if kid == "" {
		err = errors.New("kid is empty")
		return res, err
	}

	cnfObj, _ := GetConfObj()
	prefix := cnfObj.String("task_lcok_key")
	key := prefix + kid

	client, err := GetRedisClient()
	if err != nil {
		return res, err
	}

	err = client.Del(key).Err()
	if err == nil {
		res = true
	}

	return res, err
}

//浮点数比较>=
func GreaterOrEqual(a, b float64) bool {
	return math.Max(a, b) == a || math.Abs(a-b) < 0.000001
}

//执行URL任务
func RunUrlTask(tsk string, needreturn bool) (string, error) {
	var res string
	var err error

	tUrl, tPos, tPor, err := ParseTaskUrl(tsk)
	if err != nil {
		return res, err
	}

	//curl
	easy := curl.EasyInit()
	defer easy.Cleanup()

	easy.Setopt(curl.OPT_VERBOSE, 1)
    easy.Setopt(curl.OPT_FORBID_REUSE, 1)
	easy.Setopt(curl.OPT_URL, tUrl)
	easy.Setopt(curl.OPT_USERAGENT, SERNAME)
	easy.Setopt(curl.OPT_TIMEOUT, TASKMAXTIME)
	easy.Setopt(curl.OPT_TCP_KEEPALIVE, TASKMAXTIME)
	easy.Setopt(curl.OPT_PORT, tPor)
	if tPos != "" {
		easy.Setopt(curl.OPT_POST, 1)
		easy.Setopt(curl.OPT_POSTFIELDS, tPos)
		easy.Setopt(curl.OPT_POSTFIELDSIZE, len(tPos))
	}
	if needreturn {
		easy.Setopt(curl.OPT_WRITEFUNCTION, func(ptr []byte, _ interface{}) bool {
			res += string(ptr)
			return true
		})
	}

	if err = easy.Perform(); err == nil {
		res = Substr(res, 0, 1024)
	}

	return res, err
}

//执行命令行任务
func RunCmdTask(tsk string, needreturn bool) (string, error) {
	var res string
	var err error
	var out []byte

	if needreturn { //需要返回
		out, err = exec.Command("/bin/bash", "-c", tsk).CombinedOutput()
		res = Substr(string(out), 0, 1024)
	} else {
		err = exec.Command("/bin/bash", "-c", tsk).Start()
	}

	return res, err
}

//解析任务URL
func ParseTaskUrl(str string) (string, string, int, error) {
	var nUrl, nPos string
	var nPor int
	var err error
	var pd = &map[string]interface{}{}

	u, err := url.Parse(str)
	if err != nil {
		return nUrl, nPos, nPor, err
	}

	//端口号
	hs := strings.Split(u.Host, ":")
	if len(hs) == 2 {
		nPor, _ = strconv.Atoi(hs[1])
	} else if "https" == strings.ToLower(u.Scheme) {
		nPor = 443
	} else {
		nPor = 80
	}

	m, _ := url.ParseQuery(u.RawQuery)
	q := u.Query()

	for k, v := range m {
		if k == "kt_post" { //post数据
			q.Del(k)
			if v[0] != "" {
				str := strings.Trim(v[0], "\"' ")
				err = json.Unmarshal([]byte(str), pd)
				num := len(*pd)
				if err == nil && num > 0 {
					tmpU, _ := url.Parse(nPos)
					tmpQ := tmpU.Query()
					for pk, pv := range *pd {
						tmpQ.Add(pk, fmt.Sprintf("%v", pv))
					}
					nPos = tmpQ.Encode()
				} else {
					//LogRunes("json error.", "ori:", u.RawQuery, m, "str:", str, "err:", err)
				}
			}
		}
	}

	u.RawQuery = q.Encode()
	nUrl = u.String()

	return nUrl, nPos, nPor, err
}
