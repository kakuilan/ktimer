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
func CheckPidFile() (string,error) {
    var err error
    var pidfile string
    CnfObj,err = GetConfObj()
    if err !=nil {
        return "",err
    }

    pidfile = CnfObj.String("pidfile")
    if pidfile=="" {
        err = errors.New("pid path is empty.")
        return "",err
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

    return pidfile,err
}

//服务错误处理
func ServiceError(msg string, err error) {
	fmt.Println(msg, err)
	os.Exit(1)
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

func ServiceStart() {
	ServiceInit()
	TimerContainer()
}

func ServiceStop() {
	ServiceInit()

}

func ServiceRestart() {
	ServiceInit()

}

func ServiceStatus() {
	ServiceInit()

}

//版本
func ServiceVersion() {
	fmt.Printf("Version %s [%s]\n", VERSION, PUBDATE)
	os.Exit(0)
}

//服务异常处理
func ServiceException() {
	el, _ := GetErrLoger()
	if err := recover(); err != nil {
		el.Println(err)
	}

}
