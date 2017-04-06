package ktimer

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
)

//全局日志lum配置对象
var LogLumCnf = make(map[string]*lumberjack.Logger)

//获取日志的lumberjack配置
func GetLogLum(logname string) (*lumberjack.Logger, int, error) {
	var lum *lumberjack.Logger
	var open int
	var err error
	var ok bool

	CnfObj, err = GetConfObj()
	if err != nil {
		return lum, open, err
	}
	open, err = CnfObj.Int("log::log." + logname + "_open")
	if err != nil {
		open, err = 1, nil
	}

	if lum, ok = LogLumCnf[logname]; !ok {
		var logdir string
		logdir, err = CheckLogdir()
		if err != nil {
			return lum, open, err
		}

		file := CnfObj.String("log::log." + logname + "_file")
		if file == "" || err != nil {
			err, file = nil, logname+".log"
		}

		maxsize, err := CnfObj.Int("log::log.maxsize")
		if err != nil {
			maxsize, err = 500, nil
		}

		maxbackup, err := CnfObj.Int("log::log.maxbackup")
		if err != nil {
			maxbackup, err = 5, nil
		}

		maxage, err := CnfObj.Int("log::log.maxage")
		if err != nil {
			maxage, err = 30, nil
		}

		lum = &lumberjack.Logger{
			Filename:   logdir + "/" + file,
			MaxSize:    maxsize,
			MaxBackups: maxbackup,
			MaxAge:     maxage,
		}
		LogLumCnf[logname] = lum
	}

	return lum, open, err
}

//获取服务日志对象
func GetSerLoger() (*log.Logger, error) {
	var lg *log.Logger
	var err error

	CnfObj, err = GetConfObj()
	if err != nil {
		return lg, err
	}

	lum, open, err := GetLogLum("serve")
	if err != nil {
		return lg, err
	}

	if open >= 1 {
		lg = log.New(lum, "", log.Ldate|log.Lmicroseconds)
	} else {
		lg = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	}

	return lg, err
}

//获取运行日志对象
func GetRunLoger() (*log.Logger, error) {
	var lg *log.Logger
	var err error
	CnfObj, err = GetConfObj()
	if err != nil {
		return lg, err
	}

	lum, open, err := GetLogLum("runed")
	if err != nil {
		return lg, err
	}

	if open >= 1 {
		lg = log.New(lum, "", log.Ldate|log.Lmicroseconds)
	} else {
		lg = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	}

	return lg, err
}

//获取web日志对象
func GetWebLoger() (*log.Logger, error) {
	var lg *log.Logger
	var err error

	CnfObj, err = GetConfObj()
	if err != nil {
		return lg, err
	}

	lum, open, err := GetLogLum("webac")
	if err != nil {
		return lg, err
	}

	if open >= 1 {
		lg = log.New(lum, "", log.Ldate|log.Lmicroseconds)
	} else {
		lg = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	}

	return lg, err
}

//获取错误日志对象
func GetErrLoger() (*log.Logger, error) {
	var lg *log.Logger
	var err error

	CnfObj, err = GetConfObj()
	if err != nil {
		return lg, err
	}

	lum, open, err := GetLogLum("error")
	if err != nil {
		return lg, err
	}

	if open >= 1 {
		lg = log.New(lum, "", log.Ldate|log.Llongfile|log.Lmicroseconds)
	} else {
		lg = log.New(os.Stdout, "", log.Ldate|log.Llongfile|log.Lmicroseconds)
	}

	return lg, err
}

//记录服务信息日志
func LogService(v ...interface{}) {
	lg, _ := GetSerLoger()
	lg.Println(v...)
}

//记录WEB信息日志
func LogWebes(v ...interface{}) {
	lg, _ := GetWebLoger()
	lg.Println(v...)
}

//记录运行信息日志
func LogRunes(v ...interface{}) {
	lg, _ := GetRunLoger()
	lg.Println(v...)
}

//记录错误信息日志
func LogErres(v ...interface{}) {
	lg, _ := GetErrLoger()
	lg.Println(v...)
}
