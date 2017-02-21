package ktimer
import (
    "errors"
    "io/ioutil"
    "os"
    "strconv"
    "syscall"
)

//获取pid文件的值
func PidGetVue(pidfile string) (int,error) {
    value,err := ioutil.ReadFile(pidfile)
    if (err!=nil) {
        return 0,err
    }
    
    pid,err := strconv.ParseInt(string(value), 10, 32)
    if(err != nil ){
        return 0,err
    }

    return int(pid),nil
}


func PidIsActive(pid int) (bool,error) {
    if pid<=0 {
        return false,errors.New("ktimer process id error")
    }
    p,err := os.FindProcess(pid)
    if err != nil {
        return  false,err
    }

    if err := p.Signal(os.Signal(syscall.Signal(0)));err != nil{
        return false, err
    }

    return true,nil
}


func PidCreate(pidfile string) (int,error) {
    if _,err := os.Stat(pidfile);!os.IsNotExist(err) {
        if pid,_ :=  PidGetVue(pidfile);pid>0 {
            if ok,_ := PidIsActive(pid);ok {
                return pid,errors.New("ktimer pid is exists")
            } 
        }
    }

    if pf,err := os.OpenFile(pidfile, os.O_RDWR|os.O_CREATE,0600);err !=nil{
        return 0,err
    }else{
        pid := os.Getpid()
        pf.Write([]byte(strconv.Itoa(pid)))
        return pid,nil
    }
}


