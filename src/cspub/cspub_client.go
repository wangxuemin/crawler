package cspub

import (
	"baidu"
	"fmt"
	"log"
	"net"
	"os"
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

type CspubTask struct {
	Cs_task_id int64  `mcpack:"cs_task_id"`
	User       string `mcpack:"user"`
	User_key   string `mcpack:"user_key"`
	Dest_host  string `mcpack:"dest_host"`
	Dest_port  int32  `mcpack:"dest_port"`
	Target_url string `mcpack:"target_url"`
}

type CspubLogin struct {
	User_name string `mcpack:"user_name"`
	User_key  string `mcpack:"user_key"`
	Cmd_type  string `mcpack:"cmd_type"`
	Priority  int    `mcpack:"priority"`
}

type CspubClient struct {
	RemoteAddr string
	Connection net.Conn
	user_id    int32
}

func (e *CspubError) Error() string {
	return e.msg
}

func New(remoteAddr string) *CspubClient {
	c := new(CspubClient)
	c.RemoteAddr = remoteAddr

	return c
}

func (client *CspubClient) connect(retries int, timeout int) *CspubError {
	for i := 0; i < retries; i++ {
		conn, err := net.DialTimeout("tcp", client.RemoteAddr,
			time.Duration(timeout)*time.Second)
		if err != nil {
			log.Printf("connect to %s error (%d/%d): %s",
				client.RemoteAddr, i+1, retries, err.Error())
			continue
		}
		client.Connection = conn

		cslogin := CspubLogin{
			User_name: "magicnum4",
			User_key:  "magicnum498",
			Cmd_type:  "login",
			Priority:  0,
		}

		login_body, merr := baidu.Marshal(cslogin)
		if merr != nil {
			return &CspubError{"packing error: " + merr.Error()}
		}
		request := baidu.NsheadPack(login_body, 0)

		file, _ := os.Create("send")
		file.Write(request)
		file.Close()

		client.Connection.Write(request)

		response_header, response, cerr := baidu.NsheadRead(client.Connection)
		if cerr != nil {
			return &CspubError{"read pack error: " + cerr.Error()}
		}

		file, _ = os.Create("recv")
		file.Write(response_header)
		file.Write(response)
		file.Close()

		cred := CspubCred{}
		baidu.Unmarshal(response, &cred)

		task := CspubTask{
			Cs_task_id: cred.Cs_task_id,
			User:       strconv.Itoa(cred.User_id),
			User_key:   "magicnum498",
			Dest_host:  "10.210.71.14",
			Dest_port:  12345,
			Target_url: "http://open.jjwxc.net/aladdin/getPageList",
		}

		task_body, terr := baidu.Marshal(task)
		if terr != nil {
			return &CspubError{"packing error: " + terr.Error()}
		}

		grab_request := baidu.NsheadPack(task_body, 0)

		file, _ = os.Create("task_send")
		file.Write(grab_request)
		file.Close()

		client.Connection.Write(grab_request)

		return nil
	}

	log.Printf("connect to %s error after try %d times", client.RemoteAddr, retries)

	return &CspubError{"connect error"}
}

func main() {
	var client *CspubClient = New("db-spi-pubtest0.db01:7205")
	err := client.connect(3, 10)
	if err != nil {
		fmt.Println(err.Error())
	}
}
