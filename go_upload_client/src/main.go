package main

import (
	"runtime"
	"time"
	"uploadTest/src/serves"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	serves.Run()

	//ListDir()

	//	time.Sleep(5000000000)
	//task.PostUploadState(State_Pause)
	for {
		time.Sleep(100000000)
	}
}
