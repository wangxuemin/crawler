package main

import (
	"env"
	"saver"
)

func main() {
	var svr *saver.Saver
	var err error

	if err := env.Load("conf", "saver.conf"); err != nil {
		panic(err.Error())
	}
	defer env.Close()

	if svr, err = saver.NewSaver(); err != nil {
		panic("init saver error : " + err.Error())
	}

	svr.Work()
}
