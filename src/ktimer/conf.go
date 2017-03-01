package ktimer

import (
	//"github.com/astaxie/beego/config"
	//"fmt"
	"config"
	"os"
)

//默认配置
const DEFAULT_CONF = `
[default]
#pid文件
pidfile = ktimer.pid
#最大任务数量
task_max_num = 500000
#所有任务池缓存key
task_pool_key = ktimer:tasks:all
#待运转任务缓存key
task_trun_key = ktimer:tasks:second
#任务过期限制.默认执行60秒内的任务,超过则抛弃;为0则不限制,全部执行
task_expire_limit = 60

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
redis.db = 0
;redis.passwd = 

#日志
[log]
log.dir = log
log.error_open = 1
log.error_file = error.log
log.runed_open = 1
log.runed_file = runed.log
`

//全局配置对象
var CnfObj config.ConfigInterface

//获取配置文件路径
func GetConfFilePath() string {
	return GetCurrentDirectory() + "/conf.ini"
}

//检查配置文件是否存在
func CheckConfFile() bool {
	confFile := GetConfFilePath()
	res := FileExist(confFile)
	return res
}

//创建配置文件
func CreateConfFile() (bool, error) {
	res := false
	confFile := GetConfFilePath()
	fout, err := os.Create(confFile)
	defer fout.Close()
	if err == nil {
		fout.WriteString(DEFAULT_CONF)
	}

	return res, err
}

//获取配置对象
func GetConfObj() (config.ConfigInterface, error) {
	var err error
	if CnfObj == nil {
		file := GetConfFilePath()
		CnfObj, err = config.NewConfig(file)
	}

	//println("in GetConfObj", CnfObj)
	return CnfObj, err
}
