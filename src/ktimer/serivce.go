package ktimer

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/takama/daemon"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"strings"
)

//定义KT服务类型
type KTService struct {
	daemon.Daemon
}

//全局redis客户端
var SerRedisCli *redis.Client

//获取redis连接
func GetRedisClient() (*redis.Client, error) {
	var client *redis.Client
	var err error

    //当前进程是否服务进程
	isMainProc, _ := CheckCurrent2ServicePid()
	if isMainProc {
        if SerRedisCli==nil {
            SerRedisCli,err = _getRedisClient()
        }
        return SerRedisCli,err
	}else{
        client,err = _getRedisClient()
    }

	return client, err
}

func _getRedisClient() (*redis.Client, error) {
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

	defer client.Close()

	return true, err
}

//检查运行时目录
func CheckRuntimedir() (string, error) {
	var rundir string
	var err error

	currdir := GetCurrentDirectory()
	rundir = currdir + "/runtime"
	direxis := FileExist(rundir)
	if !direxis {
		err = os.MkdirAll(rundir, 0766)
		if err != nil {
			err = errors.New("failed to create runtime directiory:" + rundir)
			return rundir, err
		}
	} else {
		write := Writeable(rundir)
		err = os.Chmod(rundir, 0766)
		if !write || err != nil {
			err = errors.New("runtime dir cannot write:" + rundir)
			return rundir, err
		}
	}

	return rundir, err
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
	if pos == -1 || pos != 0 {
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

//服务退出处理
func ServiceExit(msg string, err error) {
	stacks := _getRunStack(true)
	LogErres(msg, err, string(stacks))
	if err != nil {
		fmt.Println(msg, err)
		os.Exit(1)
	} else {
		fmt.Println(msg)
		os.Exit(0)
	}
}

//服务错误处理
func ServiceError(msg string,err error) {
    stacks := _getRunStack(false)
    LogErres(msg, err, string(stacks))
}

//服务异常处理
func ServiceException() {
	if err := recover(); err != nil {
		stacks := _getRunStack(true)
		LogErres("catch panic err:", err, string(stacks))
		fmt.Println(err, string(stacks))
		//os.Exit(1)
	}
}

//获取运行栈
func _getRunStack(all bool) []byte {
	buf := make([]byte, 512)
	for {
		size := runtime.Stack(buf, all)
		if size == len(buf) {
			buf = make([]byte, len(buf)<<1)
			continue
		}
		break
	}

	return buf
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
			ServiceExit("check conf has error.", err)
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
		ServiceExit("check log`s dir has error:", err)
	}

	//检查运行时目录
	_, err = CheckRuntimedir()
	if err != nil {
		ServiceExit("check runtime dir has error:", err)
	}

	//检查pid
	pidf, err := CheckPidFile()
	if err != nil {
		ServiceExit("check pid has error:"+pidf, err)
	}

	//设置异常处理
	defer ServiceException()

}

//获取守护进程的服务对象
func GetDaemon() (*KTService, error) {
	var err error
	var dependencies = []string{"ktimer.service"}
	srv, err := daemon.New(SERNAME, SERDESC, dependencies...)
	if err != nil {
		ServiceExit("get daemon err:", err)
	}
	service := &KTService{srv}
	return service, err
}

//安装服务
func ServiceInstall() {
	ServiceInit()
	service, _ := GetDaemon()
	status, err := service.Install()
	if err != nil {
		ServiceExit("service install fail.", err)
	}
	fmt.Println(status)
	LogService("service install success.")
}

//卸载服务
func ServiceRemove() {
	ServiceInit()
	service, _ := GetDaemon()
	status, err := service.Remove()
	if err != nil {
		ServiceExit("service remove fail.", err)
	}
	fmt.Println(status)
	LogService("service remove success.")
}

//启动服务
func ServiceStart() {
	var chk bool
	var msg string

	ServiceInit()

	serPidno, _ := GetServicePidNo()
	chk, _ = PidIsActive(serPidno)
	if chk {
		msg = fmt.Sprintf("service [%d] is running,start fail.", serPidno)
		LogService(msg)
		fmt.Println(msg)
		os.Exit(0)
	}

	service, _ := GetDaemon()
	status, err := service.Start()
	if err != nil {
		ServiceExit("service start fail.", err)
	}

	//保存当前运行的配置
	rundir, _ := CheckRuntimedir()
	if curCnf, err := os.OpenFile(rundir+"/runcnf", os.O_RDWR|os.O_CREATE, 0600); err != nil {
		LogService("save run conf fail.", err)
	} else {
		cnfobj, _ := GetConfObj()
		str := fmt.Sprint(cnfobj)
		curCnf.Write([]byte(str))
	}

	fmt.Println(status)
	LogService("service start success.")
}

//停止服务
func ServiceStop() {
	ServiceInit()
	var status string
	var err error
	var chk bool

	//先检查是否在运行
	status = "service is stopped."
	serPidno, _ := GetServicePidNo()
	chk, _ = PidIsActive(serPidno)
	if serPidno == 0 || chk {
		service, _ := GetDaemon()
		status, err = service.Stop()
		if err != nil && serPidno > 0 {
			serProcess, err := os.FindProcess(serPidno)
			if err != nil {
				ServiceExit("service stop fail.", err)
			}
			if err = serProcess.Kill(); err != nil {
				ServiceExit("service process kill fail.", err)
			}
			status = "service stop success."
		}

		//删除pid
		pidfile, err := CheckPidFile()
		if FileExist(pidfile) && err == nil {
			err = os.Remove(pidfile)
			if err != nil {
				LogErres(err)
			}
		}

	}

	fmt.Println(status)
	LogService("service stop success.")
}

//查看服务状态
func ServiceStatus() {
	ServiceInit()
	service, _ := GetDaemon()
	status, err := service.Status()
	if err != nil {
		ServiceExit("service status fail.", err)
	}
	fmt.Println(status)
}

//重启服务
func ServiceRestart() {
	ServiceInit()
	LogService("service restart begining...")
	ServiceStop()
	ServiceStart()
	LogService("service restart success.")
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
	if err != nil {
		msg = fmt.Sprintf("main service create pid fail:[%s]\n", pidfile)
		ServiceExit(msg, nil)
	}
	SetCurrentServicePid(ServPidno)

	msg = fmt.Sprintf("main service run success[%d].", ServPidno)
	LogService(msg)
	fmt.Println(msg)

	//性能监控
	OpenProfile()

	TimerContainer()
	WebContainer()
}

//打开性能调试
func OpenProfile() {
	cnfobj, _ := GetConfObj()
	isOpen, _ := cnfobj.Int("profile_open")
	if isOpen >= 1 {
		port := cnfobj.String("profile_port")
		go func() {
			http.ListenAndServe("0.0.0.0:"+port, nil)
		}()
	}
}

//查看运行时服务的信息
func ServiceInfo() {
	var cnfStr string
	taskNum, _ := CountTimer()

	ServiceInit()
	serPidno, _ := GetServicePidNo()
	chk, _ := PidIsActive(serPidno)
	if !chk {
		fmt.Println("service is not running.")
		os.Exit(0)
	}

	//当前运行的配置
	runDir, _ := CheckRuntimedir()
	runCnf := runDir + "/runcnf"
	cnfBuf, err := ioutil.ReadFile(runCnf)
	if err == nil {
		cnfStr = string(cnfBuf)
	}

	fmt.Println("service is running...")
	fmt.Println("current tasks num:", taskNum)
	fmt.Println("current running conf:")
	fmt.Println(cnfStr)

}

//查看版本
func ServiceVersion() {
	fmt.Printf("Version %s [%s]\n", VERSION, PUBDATE)
	os.Exit(0)
}

//查看任务列表
func ServiceList() {

}
