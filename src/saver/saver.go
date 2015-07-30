package saver

import (
	"cspub"
	"env"
	"fmt"
	"net/rpc"
	"proto"
	"scheduler"
)

type Saver struct {
	Server          cspub.CspubResultServer
	SchedulerClient *rpc.Client

	//Config section
	Port          int    `config:"saver|port`
	SchedulerHost string `config:"schduler|host`
	SchedulerHost int    `config:"schduler|port`
}

func (*Saver) GetContextInterface() interface{} {
	return &scheduler.CrawContext{}
}

func NewSaver() (*Saver, error) {
	saver := Saver{}

	if err := env.Configure(saver); err != nil {
		return nil, err
	}
	return saver, nil
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
		go saver.SendTask(task)
	}
}

func (saver *Saver) SendTask(task *scheduler.RpcTask) {

	var reply bool
	err = client.Call("Scheduler.AddTask", task, &reply)
	if err != nil {
		fmt.Println(err)
	}
}
