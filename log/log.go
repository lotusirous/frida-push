package log

import (
	"log"
)

func Infoln(args ...interface{}) {
	a := append([]interface{}{"INFO"}, args...)
	log.Println(a...)
}
func Fatalln(args ...interface{}) {
	a := append([]interface{}{"FATAL"}, args...)
	log.Fatalln(a...)
}
