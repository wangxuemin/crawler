package main

import (
	"cspub"
	"env"
	"fmt"
	"net/rpc"
	"proto"
	"scheduler"
)

type Saver struct {
	Server cspub.CspubResultServer
}

func (*Saver) GetContextInterface() interface{} {
	return &scheduler.CrawContext{}
}

func (saver *Saver) HandleCspubResult(result *cspub.CspubFetchResult) {
	c := result.User_data.(*scheduler.CrawContext)
	switch c.Level {
	case scheduler.CRAW_LEVEL_ENTRANCE:
		saver.entrance(result)
	default:
		env.Log.Warn("unknown level %d", c.Level)
		return
	}
}

func (saver *Saver) Work() {
	saver.Server.AddResultHandler(saver)
	saver.Server.Listen()
	saver.Server.Work()
}

func (saver *Saver) entrance(result *cspub.CspubFetchResult) {
	c := result.User_data.(*scheduler.CrawContext)

	entry, err := proto.DecodeEntry(result.Html_body)
	if err != nil {
		return
	}

	for _, page := range entry.Page_list {
		page_context := scheduler.CrawContext{
			Level: scheduler.CRAW_LEVEL_PAGE,
			Cpid:  c.Cpid,
		}
		task := &scheduler.RpcTask{
			Target_url: page,
			Context:    page_context,
		}
		saver.SendTask(task)
	}
}

func (saver *Saver) SendTask(task *scheduler.RpcTask) {
	client, err := rpc.DialHTTP("tcp", "localhost:12391")
	var reply bool
	err = client.Call("scheduler.AddTask", task, &reply)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	if err := env.Load("conf", "crawler.conf"); err != nil {
		panic(err.Error())
	}

	saver := Saver{
		Server: cspub.CspubResultServer{
			Host: "10.210.71.14",
			Port: 12346,
		},
	}
	saver.Work()
}
