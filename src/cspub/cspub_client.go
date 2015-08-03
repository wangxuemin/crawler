// cspub package provides client and server to baidu cspub
package cspub

import (
	"baidu"
	"env"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"time"
)

type CspubError struct {
	msg string
}

type CspubCred struct {
	Cs_task_id int64 `mcpack:"cs_task_id"`
	User_id    int   `mcpack:"user_id"`
}

// task type to cspub server, it's be packed by mcpack struct
type CspubTask struct {
	Cs_task_id int64  `mcpack:"cs_task_id"`
	User       string `mcpack:"user"`
	Userkey    string `mcpack:"user_key"`
	Dest_host  string `mcpack:"dest_host"`
	Dest_port  int32  `mcpack:"dest_port"`

	//field belong should be set by caller
	Target_url string      `mcpack:"target_url"`
	Priority   int32       `mcpack:"priority"`
	User_data  interface{} `mcpack:"trespassing_field"` //User defined data
}

type CspubLogin struct {
	Username string `mcpack:"user_name"`
	Userkey  string `mcpack:"user_key"`
	Cmd_type string `mcpack:"cmd_type"`
	Priority int    `mcpack:"priority"`
}

type CspubClient struct {
	Cred       CspubCred
	Connection net.Conn

	Username  string
	Userkey   string
	Retry     int
	Timeout   int
	Dest_host string
	Dest_port int32
	User_id   int32
}

const (
	CMD_LOGIN = "login"
)

func (e *CspubError) Error() string {
	return e.msg
}

func (client *CspubClient) Connect(addr string) *CspubError {
	for i := 0; i < client.Retry; i++ {
		conn, err := net.DialTimeout("tcp", addr,
			time.Duration(client.Timeout)*time.Second)
		if err != nil {
			log.Printf("connect to %s error (%d/%d): %s",
				addr, i+1, client.Retry, err.Error())
			continue
		}
		client.Connection = conn

		cslogin := CspubLogin{
			Username: client.Username,
			Userkey:  client.Userkey,
			Cmd_type: CMD_LOGIN,
			Priority: 0,
		}

		login_body, merr := baidu.Marshal(cslogin)
		if merr != nil {
			return &CspubError{"packing error: " + merr.Error()}
		}
		request := baidu.NsheadPack(login_body, 0)
		client.Connection.Write(request)
		_, response, cerr := baidu.NsheadRead(client.Connection)
		if cerr != nil {
			return &CspubError{"read pack error: " + cerr.Error()}
		}

		baidu.Unmarshal(response, &client.Cred)
		return nil
	}

	env.Log.Warn("connect to %s error after try %d times", addr, client.Retry)

	return &CspubError{"connect error"}
}

func (client *CspubClient) SendTask(task CspubTask) *CspubError {
	task.Cs_task_id = client.Cred.Cs_task_id
	task.User = strconv.Itoa(client.Cred.User_id)
	task.Userkey = client.Userkey
	task.Dest_host = client.Dest_host
	task.Dest_port = client.Dest_port

	task_body, terr := baidu.Marshal(task)
	if terr != nil {
		env.Log.Info("packing error %s", terr.Error())
		return &CspubError{"packing error: " + terr.Error()}
	}

	grab_request := baidu.NsheadPack(task_body, 0)
	ioutil.WriteFile("cspub_request", grab_request, 0644)
	client.Connection.Write(grab_request)
	env.Log.Info("sended task to cspub done [url:%s] [priority:%d]",
		task.Target_url, task.Priority)

	return nil
}
