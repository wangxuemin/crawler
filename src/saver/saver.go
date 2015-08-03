package saver

import (
	"cspub"
	"env"
	"fmt"
	"net/rpc"
	"scheduler"
	"time"
)

type Saver struct {
	Server          cspub.CspubResultServer
	SchedulerClient *rpc.Client

	//Config section
	Port          int    `config:"saver|port"`
	SchedulerHost string `config:"scheduler|host"`
	SchedulerPort int    `config:"scheduler|port"`
}

func (*Saver) GetContextInterface() interface{} {
	return &scheduler.CrawContext{}
}

func NewSaver() (*Saver, error) {
	saver := &Saver{}

	var err error
	if err = env.Configure(saver); err != nil {
		return nil, err
	}

	saver.Server.AddResultHandler(saver)
	if err = saver.Server.Listen(saver.Port); err != nil {
		env.Log.Critical("bind on port %d error", saver.Port)
		return nil, err
	}

	return saver, nil
}

func (saver *Saver) HandleCspubResult(result *cspub.CspubFetchResult) {
	c := result.User_data.(*scheduler.CrawContext)

	switch c.Level {
	case scheduler.CRAW_LEVEL_ENTRANCE:
		saver.entrance(result)
	case scheduler.CRAW_LEVEL_PAGE:
		saver.page(result)
	case scheduler.CRAW_LEVEL_NOVEL:
		saver.novel(result)
	default:
		env.Log.Warn("unknown level %d", c.Level)
		return
	}
}

func (saver *Saver) Work() {
	saver.Server.Work()
}

func (saver *Saver) SendTask(task *scheduler.RpcTask) {
	var reply bool

	if saver.SchedulerClient == nil {
		saver.redail()
	}

	for {
		err := saver.SchedulerClient.Call("Scheduler.AddTask", task, &reply)
		if err != nil {
			saver.redail() // TODO this will cause all saver to call redail
			env.Log.Critical("error: " + err.Error())
		}

		env.Log.Info("sended to scheduler task [%s] [%t]", task.Target_url, reply)
		break
	}
}

func (saver *Saver) redail() {
	for {
		var err error
		if saver.SchedulerClient, err = rpc.DialHTTP("tcp",
			fmt.Sprintf("%s:%d", saver.SchedulerHost, saver.SchedulerPort),
		); err != nil {
			env.Log.Warn("dial to scheduler error: " + err.Error())
			time.Sleep(time.Second * 5)
		} else {
			break
		}
	}
}
