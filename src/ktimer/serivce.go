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
func CheckLogdir() (bool, error) {
	var err error
	CnfObj, err = GetConfObj()
	if err != nil {
		return false, err
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
			return false, err
		}
	} else {
		write := Writeable(logdir)
		err = os.Chmod(logdir, 0766)
		if !write || err != nil {
			err = errors.New("logdir canot write:" + logdir)
			return false, err
		}
	}

	return true, err
}

//检查pid文件
func CheckPid() (bool, error) {
	var err error
	var chk bool = false
	CnfObj, err = GetConfObj()

	pid := CnfObj.String("pidfile")
	if pid == "" {
		err = errors.New("pid path is empty.")
		return chk, err
	}
	pid = strings.Replace(pid, "\\", "/", -1)
	pos := strings.Index(pid, "/")
	if pos == -1 {
		currdir := GetCurrentDirectory()
		pid = currdir + "/" + strings.TrimRight(pid, "/")
	}

	if !FileExist(pid) {
		_, err = PidCreate(pid)
		if err != nil {
			err = errors.New("failed to create pid file:" + pid)
			return chk, err
		}
	}

	return chk, err
}

//服务错误处理
func ServiceError(msg string, err error) {
	fmt.Println(msg, err)
	os.Exit(0)
}

//初始化(第一次执行时)
func ServiceInit() {
	var err error
	var chk bool

	fmt.Println("CnfObj", CnfObj)

	//检查配置文件
	chk = CheckConfFile()
	if !chk {
		CreateConfFile()
	}

	//检查redis
	chk, err = CheckRedis()
	if err != nil {
		ServiceError("redis connet has error:", err)
	}

	//检查日志目录
	chk, err = CheckLogdir()
	if err != nil {
		ServiceError("check log`s dir has error:`", err)
	}

	//检查pid
	chk, err = CheckPid()
	if err != nil {
		ServiceError("check pid has error:", err)
	}

}

func ServiceStart() {
	ServiceInit()

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

func ServiceVersion() {
	fmt.Printf("Version %s [%s]\n", VERSION, PUBDATE)
	os.Exit(0)
}
