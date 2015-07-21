package main

import (
	//"cspub"
	"github.com/alecthomas/log4go"
	//"fmt"
)

type Job struct {
}

/*
func (*Job) HandleCspubResult(result *cspub.CspubFetchResult) {
	db, err = sql.Open("mysql", "root:root@localhost/novels?charset=utf8")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(result.Html_body))
}
*/

func main() {
	if err := Load("conf", "crawler.conf"); err != nil {
		panic(err.Error())
	}
	defer Close()

	log4go.Info("Hello")
	log4go.Warn("Oops")
	/*
		server := cspub.CspubResultServer{}
		server.Listen(":12345")

		server.AddResultHandler(&Job{})

		go server.Work()

		client := cspub.CspubClient{}
		username := "magicnum4"
		user_key := "magicnum498"
		err := client.Connect("db-spi-pubtest0.db01:7205", username, user_key, 3, 10)
		if err != nil {
			fmt.Println(err.Error())
		}
	*/

}
