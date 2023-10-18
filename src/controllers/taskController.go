package controllers

import (
	"fmt"
	"strconv"
	"time"

	middlewares "orlangur.link/services/mini.note/handlers"
)

// Tasker -> run task after timeout
func Tasker() {
	importInterval, _ := strconv.Atoi(middlewares.DotEnvVariable("UPDATE_IMPORT_HOUR", "1"))
	importTicker := time.NewTicker(time.Duration(importInterval) * time.Hour)
	importMinute := time.NewTicker(time.Minute)
	for {
		select {
		case <-importTicker.C:
			fmt.Println("task")
		case <-importMinute.C:
			fmt.Println("minute")
		}
	}
}
