package handler

import "log"

var Logger *log.Logger

func logPrintln(v ...interface{}) {
	if Logger == nil {
		log.Println(v...)
		return
	}
	Logger.Println(v...)
}
