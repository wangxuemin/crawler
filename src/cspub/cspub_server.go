package cspub

import (
	"C"
	"baidu"
	"env"
	"fmt"
	"net"
)

type CspubFetchResult struct {
	Target_url     string      `mcpack:"target_url"`
	Result         int32       `mcpack:"result"`
	User           string      `mcpack:"user"`
	Status         string      `mcpack:"status"`
	Method         string      `mcpack:"method"`
	Cur_url        string      `mcpack:"cur_url"`
	Html_body      []byte      `mcpack:"html_body"`
	Trunc_overflow int         `mcpack:"trunc_overflow_page"`
	User_data      interface{} `mcpack:"trespassing_field"`
}

type CspubResultReceiver interface {
	GetContextInterface() interface{}
	HandleCspubResult(result *CspubFetchResult)
}

type CspubResultServer struct {
	listener  net.Listener
	receivers []CspubResultReceiver
}

func (server *CspubResultServer) Listen(port int) error {
	if listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port)); err != nil {
		return err
	} else {
		server.listener = listener
	}

	return nil
}

func (server *CspubResultServer) AddResultHandler(receiver CspubResultReceiver) {
	server.receivers = append(server.receivers, receiver)
}

func (server *CspubResultServer) Work() error {
	for {
		conn, err := server.listener.Accept()
		env.Log.Debug("new Connection Accept [%s]", conn.RemoteAddr().String())
		if err != nil {
			env.Log.Critical("listen error " + err.Error())
			break
		}
		channel := make(chan string)
		go server.request_handler(conn, channel)
	}

	env.Log.Info("Server Exit")
	return nil
}

func (server *CspubResultServer) request_handler(conn net.Conn, out chan string) {
	for {
		_, resp, err := baidu.NsheadRead(conn)
		if err != nil {
			env.Log.Warn("read error " + err.Error())
			break
		}

		for _, receiver := range server.receivers {
			context := receiver.GetContextInterface()
			fetchResult := &CspubFetchResult{
				User_data: context,
			}
			if err := baidu.Unmarshal(resp, fetchResult); err != nil {
				env.Log.Warn("unmarshal resp %s error", string(resp))
				break
			}

			go receiver.HandleCspubResult(fetchResult)
		}
	}
}
