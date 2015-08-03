package main

import (
	"env"
	"scheduler"
)

func main() {
	if err := env.Load("conf", "scheduler.conf"); err != nil {
		panic(err.Error())
	}
	defer env.Close()

	sch, err := scheduler.NewScheduler()
	if err != nil {
		panic("init scheduler error: " + err.Error())
	}

	sch.Work()

}
