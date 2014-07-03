package main

import (
	"fmt"
	"log"
)

func Info(v ...interface{}) {
	log.Println("[INFO]", fmt.Sprint(v...))
}

func Warn(v ...interface{}) {
	log.Println("[WARN]", fmt.Sprint(v...))
}

func Fatal(v ...interface{}) {
	log.Fatal("[FATAL]", fmt.Sprint(v...))
}
