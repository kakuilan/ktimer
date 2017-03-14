package ktimer

import (
	"fmt"
	//"github.com/astaxie/beego/config"
	//"config"
	"errors"
	"gopkg.in/redis.v5"
	"os"
	"strings"
    "github.com/takama/daemon"
)

//定义KT服务类型
type KTService struct {
    daemon.Daemon
}

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

//检查redis是否可连接
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
    el, _ := GetErrLoger()
    el.Println(msg,err)
    if err != nil {
		fmt.Println(msg, err)
		os.Exit(1)
	} else {
		fmt.Println(msg)
		os.Exit(0)
	}
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

//获取守护进程的服务对象
func GetDaemon() (*KTService,error) {
    var err error
    var dependencies = []string{"ktimer.service"}
    srv,err := daemon.New(SERNAME, SERDESC, dependencies...)
    if err!=nil {
        ServiceError("get daemon err:", err)
    }
    service := &KTService{srv}
    return service,err
} 

//安装服务
func ServiceInstall() {
    ServiceInit()
    service,_ := GetDaemon()
    status, err := service.Install()
    if err != nil {
        ServiceError("service install fail.",err)
    }
    fmt.Println(status)
}

//卸载服务
func ServiceRemove(){
   ServiceInit()
    service,_ := GetDaemon()
    status,err := service.Remove()
    if err != nil {
        ServiceError("service remove fail.",err)
    }
    fmt.Println(status)
}

//启动服务
func ServiceStart() {
    ServiceInit()
    service,_ := GetDaemon()
    status, err := service.Start()
    if err != nil {
        ServiceError("service start fail.",err)
    }
    fmt.Println(status)
}

//停止服务
func ServiceStop() {
    ServiceInit()
    service,_ := GetDaemon()
    status, err := service.Stop()
    if err != nil {
        ServiceError("service stop fail.",err)
    }else{
        //删除pid
        pidfile, err := CheckPidFile()
        if err==nil {
            err = os.Remove(pidfile)
            if err!=nil {
                ServiceError("pid file remove error.", err)
            }
        }
    }
    fmt.Println(status)
}

//查看服务状态
func ServiceStatus() {
    ServiceInit()
    service,_ := GetDaemon()
    status, err := service.Status()
    if err != nil {
        ServiceError("service status fail.",err)
    }
    fmt.Println(status)
}

//重启服务
func ServiceRestart() {
    ServiceStop()
    ServiceStart()
}

//主体服务
func ServiceMain() {
    var chk bool
    var err error
    var msg string
    
    ServiceInit()
    
    //检查pid
    chk, err = CheckCurrent2ServicePid()
    servpidno, _ := GetServicePidNo()
    servIsRun, _ := PidIsActive(servpidno)
    if chk || servIsRun {
        fmt.Printf("ktimer service [%d] is running...\n", servpidno)
        os.Exit(0)
    }

    pidfile, _ := CheckPidFile()
    ServPidno, err = PidCreate(pidfile)
    if err!= nil {
        msg = fmt.Sprintf("main service create pid fail:[%s]\n", pidfile)
        ServiceError(msg, nil)
    }
    SetCurrentServicePid(ServPidno)
   
    msg = fmt.Sprintf("main service run success[%d].", ServPidno)
    rlg, _ := GetRunLoger()
    fmt.Println(msg)
    rlg.Println(msg)
    TimerContainer()
    WebContainer() 
}

//查看运行时服务的信息
func ServiceInfo() {
    
}

//查看版本
func ServiceVersion() {
	fmt.Printf("Version %s [%s]\n", VERSION, PUBDATE)
	os.Exit(0)
}

