package ktimer

import (
	"fmt"
	//"github.com/astaxie/beego/config"
	//"config"
	"errors"
	"gopkg.in/redis.v5"
	"os"
	"strings"
)

//获取redis连接
func GetRedisClient() (*redis.Client, error) {
	var client *redis.Client
	var err error
	CnfObj, err = GetConfObj()

	if err == nil {
		var addr string
		host := CnfObj.String("redis::redis.host")
		port := CnfObj.String("redis::redis.port")
		pawd := CnfObj.String("redis::redis.passwd")
		db, err := CnfObj.Int("redis::redis.db")
		addr = host + ":" + port
		//fmt.Println(host, port, addr, pawd, db, err2, "redis conf")
		if err != nil {
			err = errors.New("read config failed,key [redis::redis.db].")
			return client, err
		}

		client = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: pawd,
			DB:       db,
		})

		return client, err
	}

	return client, err
}

//检查redis是否连接
func CheckRedis() (bool, error) {
	var client *redis.Client
	var err error
	var pong string
	var res bool = false

	client, err = GetRedisClient()
	if err != nil {
		err = errors.New("failed to get get redis client connection.")
		return res, err
	}

	pong, err = client.Ping().Result()
	if err != nil {
		return res, err
	} else if pong != "PONG" {
		err = errors.New("reids ping result not eq `PONG`.")
		return res, err
	}

	return true, err
}

//检查日志目录
func CheckLogdir() (string, error) {
	var err error
	CnfObj, err = GetConfObj()
	if err != nil {
		return "", err
	}

	logdir := CnfObj.String("log::log.dir")
	logdir = strings.Replace(logdir, "\\", "/", -1)
	pos := strings.Index(logdir, "/")
	if pos == -1 { //相对当前目录
		currdir := GetCurrentDirectory()
		logdir = currdir + "/" + strings.TrimRight(logdir, "/")
	}

	direxis := FileExist(logdir)
	if !direxis {
		err = os.MkdirAll(logdir, 0766)
		if err != nil {
			err = errors.New("failed to create log directory:" + logdir)
			return "", err
		}
	} else {
		write := Writeable(logdir)
		err = os.Chmod(logdir, 0766)
		if !write || err != nil {
			err = errors.New("logdir canot write:" + logdir)
			return "", err
		}
	}

	return logdir, err
}

//检查pid文件并返回路径
func CheckPidFile() (string, error) {
	var err error
	var pidfile string
	CnfObj, err = GetConfObj()
	if err != nil {
		return "", err
	}

	pidfile = CnfObj.String("pidfile")
	if pidfile == "" {
		err = errors.New("pid path is empty.")
		return "", err
	}

	pidfile = strings.Replace(pidfile, "\\", "/", -1)
	pos := strings.Index(pidfile, "/")
	if pos == -1 {
		currdir := GetCurrentDirectory()
		pidfile = currdir + "/" + strings.TrimRight(pidfile, "/")
	}

	piddir := GetParentDirectory(pidfile)
	chk := Writeable(piddir)
	if !chk {
		err = errors.New("pid`dir cannot be written:" + piddir)
	}

	return pidfile, err
}

//服务错误处理
func ServiceError(msg string, err error) {
	if err != nil {
		fmt.Println(msg, err)
		os.Exit(1)
	} else {
		fmt.Println(msg)
		os.Exit(0)
	}
}

//初始化检查
func ServiceInit() {
	var err error
	var chk bool

	//检查配置文件
	chk = CheckConfFile()
	if !chk {
		chk, err = CreateConfFile()
		if err != nil {
			conf := GetConfFilePath()
			err = errors.New("conf file does not exist,and create failed:" + conf)
			ServiceError("check conf has error.", err)
		}
	}

	//检查redis
	chk, err = CheckRedis()
	if err != nil {
		ServiceError("redis connet has error:", err)
	}

	//检查日志目录
	_, err = CheckLogdir()
	if err != nil {
		ServiceError("check log`s dir has error:`", err)
	}

	//检查pid
	_, err = CheckPidFile()
	if err != nil {
		ServiceError("check pid has error:", err)
	}

	fmt.Println("CnfObj", CnfObj)
}

//安装服务
func ServiceInstall() {
    
}

//卸载服务
func ServiceRemove(){
    
}

//启动服务
func ServiceStart() {
	var chk bool
	var err error
	ServiceInit()
	//检查pid
	chk, err = CheckCurrent2ServicePid()
	if chk {
		ServiceError("current process and service are the same,start fail.", nil)
	}

	ServPidno, _ := GetServicePidNo()
	servIsRun, _ := PidIsActive(ServPidno)
	if servIsRun {
		ServiceError("service is running,start fail.", nil)
	}

	pidfile, _ := CheckPidFile()
	ServPidno, err = PidCreate(pidfile)
	if err != nil {
		ServiceError("failed to create file during service startup.", nil)
	}
	SetCurrentServicePid(ServPidno)

	msg := fmt.Sprintf("service [%d] start success.", ServPidno)
	rl, _ := GetRunLoger()
	fmt.Println(msg)
	rl.Println(msg)

	TimerContainer()
}

//停止服务
func ServiceStop() {
	var err error
	ServiceInit()

	ServPidno, _ := GetServicePidNo()
	servIsRun, _ := PidIsActive(ServPidno)
	if !servIsRun {
		ServiceError("service not running.", nil)
	}

	//停止服务进程
	serProcess, err := os.FindProcess(ServPidno)
	if err != nil {
		ServiceError("service process cannot find.", err)
	}
	//if err = serProcess.Release();err!=nil {
	//ServiceError("service process release fail.", err)
	//}
	if err = serProcess.Kill(); err != nil {
		ServiceError("service process kill fail.", err)
	}

	//删除pid
	pidfile, err := CheckPidFile()
	if err != nil {
		ServiceError("check pif file has error.", err)
	}

	err = os.Remove(pidfile)
	if err != nil {
		ServiceError("pid file remove error.", err)
	}

	msg := fmt.Sprintf("service [%d] stop success.", ServPidno)
	rl, _ := GetRunLoger()
	fmt.Println(msg)
	rl.Println(msg)
	os.Exit(0)
}

//重启服务
func ServiceRestart() {
	pid := os.Getpid()
	msg := fmt.Sprintf("service restarting... currentPid[%d]", pid)
	rl, _ := GetRunLoger()
	fmt.Println(msg)
	rl.Println(msg)

	ServPidno, _ := GetServicePidNo()
	servIsRun, _ := PidIsActive(ServPidno)
	if servIsRun {
		ServiceStop()
	} else {
		msg = "service not running."
		fmt.Printf(msg)
		rl.Println(msg)
	}

	ServiceStart()
}

//查看服务状态
func ServiceStatus() {
	ServiceInit()

	ServPidno, _ := GetServicePidNo()
	servIsRun, _ := PidIsActive(ServPidno)
	if servIsRun {
		fmt.Printf("service [%d] is running.\n", ServPidno)
	} else {
		fmt.Println("service is not running.")
	}

	os.Exit(0)
}

//查看版本
func ServiceVersion() {
	fmt.Printf("Version %s [%s]\n", VERSION, PUBDATE)
	os.Exit(0)
}

//查看运行时服务的信息
func ServiceInfo() {
    
}

//服务异常处理
func ServiceException() {
	el, _ := GetErrLoger()
	if err := recover(); err != nil {
		fmt.Println(err)
		el.Println(err)
		os.Exit(1)
	}
}
