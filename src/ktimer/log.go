package ktimer

import (
	//"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	//"reflect"
)

//全局日志对象
var RunLoger, ErrLoger *lumberjack.Logger

//获取运行日志对象
func GetRunLoger() (interface{}, error) {
	var err error
	if RunLoger == nil {
		CnfObj, err = GetConfObj()
		if err != nil {
			return RunLoger, err
		}
		var logdir string
		logdir, err = CheckLogdir()
		if err != nil {
			return RunLoger, err
		}

		file := CnfObj.String("log::log.runed_file")
		maxsize, err := CnfObj.Int("log::log.maxsize")
		if err != nil {
			maxsize = 500
		}
		maxbackup, err := CnfObj.Int("log::log.maxbackup")
		if err != nil {
			maxbackup = 5
		}
		maxage, err := CnfObj.Int("log::log.maxage")
		if err != nil {
			maxage = 30
		}

		RunLoger = &lumberjack.Logger{
			Filename:   logdir + "/" + file,
			MaxSize:    maxsize,
			MaxBackups: maxbackup,
			MaxAge:     maxage,
		}

	}

	return RunLoger, err
}

//获取错误日志对象
func GetErrLoger() (interface{}, error) {
	var err error
	if ErrLoger == nil {
		CnfObj, err = GetConfObj()
		if err != nil {
			return ErrLoger, err
		}
		var logdir string
		logdir, err = CheckLogdir()
		if err != nil {
			return ErrLoger, err
		}

		file := CnfObj.String("log::log.error_file")
		maxsize, err := CnfObj.Int("log::log.maxsize")
		if err != nil {
			maxsize = 500
		}
		maxbackup, err := CnfObj.Int("log::log.maxbackup")
		if err != nil {
			maxbackup = 5
		}
		maxage, err := CnfObj.Int("log::log.maxage")
		if err != nil {
			maxage = 30
		}

		ErrLoger = &lumberjack.Logger{
			Filename:   logdir + "/" + file,
			MaxSize:    maxsize,
			MaxBackups: maxbackup,
			MaxAge:     maxage,
		}

	}

	return ErrLoger, err
}
