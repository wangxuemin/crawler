package main

import (
	"env"
	"scheduler"
)

func main() {
	if err := env.Load("conf", "crawler.conf"); err != nil {
		panic(err.Error())
	}

	sch, _ := scheduler.NewScheduler()

	sch.Work()
}
