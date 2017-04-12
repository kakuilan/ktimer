package ktimer

import (
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"syscall"
)

//全局pid变量
var ServPidno int

//获取pid文件的值
func PidGetVue(pidfile string) (int, error) {
	value, err := ioutil.ReadFile(pidfile)
	if err != nil {
		return 0, err
	}

	pid, err := strconv.ParseInt(string(value), 10, 32)
	if err != nil {
		return 0, err
	}

	return int(pid), nil
}

//检查pid进程是否存在
func PidIsActive(pid int) (bool, error) {
	if pid <= 0 {
		return false, errors.New("ktimer process id error")
	}
	p, err := os.FindProcess(pid)
	if err != nil {
		return false, err
	}

	if err := p.Signal(os.Signal(syscall.Signal(0))); err != nil {
		return false, err
	}

	return true, nil
}

//创建pid文件
func PidCreate(pidfile string) (int, error) {
	if _, err := os.Stat(pidfile); !os.IsNotExist(err) {
		if pid, _ := PidGetVue(pidfile); pid > 0 {
			if ok, _ := PidIsActive(pid); ok {
				return pid, errors.New("ktimer pid is exists")
			}
		}
	}

	if pf, err := os.OpenFile(pidfile, os.O_RDWR|os.O_CREATE, 0600); err != nil {
		return 0, err
	} else {
		pid := os.Getpid()
		pf.Write([]byte(strconv.Itoa(pid)))
		return pid, nil
	}
}

//获取服务的pid进程号
func GetServicePidNo() (int, error) {
	var pidno int
	var err error
	var pidfile string

	pidfile, err = CheckPidFile()
	if err != nil {
		return 0, err
	}

	pidno, err = PidGetVue(pidfile)
	if err != nil {
		return 0, err
	}

	_, err = PidIsActive(pidno)
	if err != nil {
		return 0, err
	}

	return pidno, err
}

//比较当前进程和服务进程的pid
func CheckCurrent2ServicePid() (bool, error) {
	var chk bool
	var err error
	var serPid, curPid int

	serPid, err = GetServicePidNo()
	if err != nil {
		return false, err
	}

	curPid = os.Getpid()
	if serPid == curPid {
		chk = true
	}

	return chk, err
}

//设置当前服务pid
func SetCurrentServicePid(pid int) {
	ServPidno = pid
}

//获取当前服务pid
func GetCurrentServicePid() int {
	return ServPidno
}
