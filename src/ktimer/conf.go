package ktimer

import (
	//"github.com/astaxie/beego/config"
	//"fmt"
	"config"
	"errors"
	"os"
    "time"
    "math/rand"
    "regexp"
)

//默认配置
const DEFAULT_CONF = `
[default]
#pid文件
pidfile = runtime/ktimer.pid
#最大任务数量
task_max_num = 500000
#最大任务天数
task_max_day = 100
#所有任务池缓存key
task_pool_key = ktimer:tasks:all
#待运转任务缓存key
task_trun_key = ktimer:tasks:second
#任务锁key
task_lcok_key = ktimer:tasks:lock
#任务过期限制.默认执行60秒内的任务,超过则抛弃;为0则不限制,全部执行
task_expire_limit = 60
#相同定时器的间隔限制.默认10秒内,若有多个相同定时器,只保留最后那个.为0则不限制.
task_sametimer_interval = 10
#相同断续器的间隔限制.默认0为不允许存在多个相同断续器.
task_sameticker_interval = 0

[web]
#是否启用web
web.enable = 1
#web端口
web.port = 9558
#web监听IP
web.bind_ip = 127.0.0.1
#web访问密码
web.passwd = eCvN5BxH$bJc

[redis]
redis.host = 127.0.0.1
redis.port = 6379
redis.db = 1
redis.passwd = 

[log]
log.dir = log
#日志最大尺寸M
log.maxsize = 500
#日志最多备份
log.maxbackup = 5
#日志保留天数
log.maxage = 30
#服务日志
log.serve_open = 1
log.serve_file = serve.log
#错误日志
log.error_open = 1
log.error_file = error.log
#运行日志
log.runed_open = 1
log.runed_file = runed.log
#web日志
log.webac_open = 1
log.webac_file = webac.log
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
        //生成新的web密码
        pwdbyt := RandPwd(12,3)
        pwdstr := "web.passwd = " + string(pwdbyt) +"\n"
        re,_ := regexp.Compile(`web.passwd( )?=( )?.*\n`)
        pwdstr = re.ReplaceAllString(DEFAULT_CONF, pwdstr)
        fout.WriteString(pwdstr)
	}

	return res, err
}

//获取配置对象
func GetConfObj() (config.ConfigInterface, error) {
	var err error
	if CnfObj == nil {
		file := GetConfFilePath()
		CnfObj, err = config.NewConfig(file)
		if err != nil {
			err = errors.New("failed to read config file:" + file)
		}
	}

	return CnfObj, err
}

//生成随机密码
func RandPwd(size int, kind int) []byte {
    //kind:0纯数字,1小写字母,2大写字母,3数字和大小写
    ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if is_all { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}
