package controllers

import (
	"fmt"
	"strconv"
	"time"

	middlewares "orlangur.link/services/mini.note/handlers"
)

//Tasker -> run task after timeout
func Tasker() {
	importInterval, _ := strconv.Atoi(middlewares.DotEnvVariable("UPDATE_IMPORT_HOUR"))
	importTicker := time.NewTicker(time.Duration(importInterval) * time.Hour)
	for {
		select {
		case <-importTicker.C:
			fmt.Println("task")
		}
	}
}
