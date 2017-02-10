package ktimer
import (
    //"github.com/astaxie/beego/config"
    //"fmt"
    "os"
)

const DEFAULT_CONF = `
[default]
#pid文件
pidfile = /var/run/ktimer.pid
#最大任务数量
task_max_num = 100000
#任务分隔符,形如 类型@时间@任务@次数限制
#例如定时器(定时器无次数限制,仅一次) timer@30@echo -e Hello Ktimer
#例如断续器(次数限制为0时即不限制) ticker@10@date --rfc-3339=ns@0
task_separator = @
#所有任务池缓存key
task_pool_key = ktimer:tasks:all
#待运转任务缓存key
task_trun_key = ktimer:tasks:second

[web]
#是否启用web
web_enable = 1
#web端口
web_port = 9558
#web监听IP
web_bind_ip = 127.0.0.1
#web访问密码
web_passwd = 123456

[redis]
redis.host = 127.0.0.1
redis.port = 6379
;redis.passwd = 

#默认任务
[task]
default_timer = timer@1@echo -e Hello Ktimer
default_ticker = ticker@1@date --rfc-3339=ns@0

#日志
[log]
log.dir = log
log.error_open = 1
log.error_file = error.log
log.runed_open = 1
log.runed_file = runed.log
`

func GetConfFilePath() string {
    return GetCurrentDirectory() + "/conf.ini"
}

func CheckConfFile() bool {
    confFile := GetConfFilePath()
    res := FileExist(confFile)
    return res
}

func CreateConfFile() (bool,error) {
    res := false
    confFile := GetConfFilePath()
    fout,err := os.Create(confFile)
    defer fout.Close()
    if(err ==nil) {
        fout.WriteString(DEFAULT_CONF)
    }

    return res,err
}
