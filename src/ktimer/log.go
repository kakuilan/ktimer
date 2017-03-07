package ktimer

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
    //"fmt"
)

//全局日志对象
var RunLoger, ErrLoger *log.Logger

//获取lumberjack运行日志配置
func GetRunLum() (*lumberjack.Logger, error) {
	var err error
	var lum *lumberjack.Logger
	CnfObj, err = GetConfObj()
	if err != nil {
		return lum, err
	}
	var logdir string
	logdir, err = CheckLogdir()
	if err != nil {
		return lum, err
	}

	file := CnfObj.String("log::log.runed_file")
	maxsize, err := CnfObj.Int("log::log.maxsize")
	if err != nil {
		maxsize,err = 500,nil
	}
	maxbackup, err := CnfObj.Int("log::log.maxbackup")
	if err != nil {
		maxbackup,err = 5,nil
	}
	maxage, err := CnfObj.Int("log::log.maxage")
	if err != nil {
		maxage,err = 30,nil
	}

	lum = &lumberjack.Logger{
		Filename:   logdir + "/" + file,
		MaxSize:    maxsize,
		MaxBackups: maxbackup,
		MaxAge:     maxage,
	}

	return lum, err
}

//获取lumberjack错误日志配置
func GetErrLum() (*lumberjack.Logger, error) {
	var err error
	var lum *lumberjack.Logger
	CnfObj, err = GetConfObj()
	if err != nil {
		return lum, err
	}
	var logdir string
	logdir, err = CheckLogdir()
	if err != nil {
		return lum, err
	}

	file := CnfObj.String("log::log.error_file")
	maxsize, err := CnfObj.Int("log::log.maxsize")
	if err != nil {
		maxsize,err = 500,nil
	}
	maxbackup, err := CnfObj.Int("log::log.maxbackup")
	if err != nil {
		maxbackup,err = 5,nil
	}
	maxage, err := CnfObj.Int("log::log.maxage")
	if err != nil {
		maxage,err = 30,nil
	}

	lum = &lumberjack.Logger{
		Filename:   logdir + "/" + file,
		MaxSize:    maxsize,
		MaxBackups: maxbackup,
		MaxAge:     maxage,
	}

	return lum, err
}

//获取运行日志对象
func GetRunLoger() (*log.Logger,error) {
    var err error
    var l *log.Logger
    lum,err := GetRunLum()
    if err!=nil {
        return l,err
    }
    l = log.New(lum, "", log.Ldate|log.Lmicroseconds)
    return l,err
}

//获取错误日志对象
func GetErrLoger() (*log.Logger,error) {
    var err error
    var l *log.Logger
    lum,err := GetErrLum()
    if err != nil {
        return l,err
    }
    l = log.New(lum, "", log.Ldate|log.Lmicroseconds)
    return l,err
}

