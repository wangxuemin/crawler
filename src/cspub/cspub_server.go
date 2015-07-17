package cspub

import (
	"C"
	"baidu"
	"log"
	"net"
)

type CspubFetchResult struct {
	Target_url string `mcpack:"target_url"`
	Result     int32  `mcpack:"result"`
	User       string `mcpack:"user"`
	Status     string `mcpack:"status"`
	Method     string `mcpack:"method"`
	Cur_url    string `mcpack:"cur_url"`
	Html_body  []byte `mcpack:"html_body"`
}

type CspubResultReciver interface {
	HandleCspubResult(result *CspubFetchResult) error
}

type CspubResultServer struct {
	listener net.Listener
	recivers []CspubResultReciver
}

func (server *CspubResultServer) Listen(addr string) error {
	if listener, err := net.Listen("tcp", addr); err != nil {
		return err
	} else {
		server.listener = listener
	}

	return nil
}

func (server *CspubResultServer) Work() {
	for {
		conn, err := server.listener.Accept()
		if err != nil {
			log.Fatal("listen error " + err.Error())
		}
		channel := make(chan string)
		go server.request_handler(conn, channel)
	}
}

func (server *CspubResultServer) request_handler(conn net.Conn, out chan string) {
	_, resp, err := baidu.NsheadRead(conn)
	if err != nil {
		log.Println("read error " + err.Error())
		return
	}

	fetchResult := &CspubFetchResult{}
	baidu.Unmarshal(resp, fetchResult)

	for _, reciver := range server.recivers {
		go reciver.HandleCspubResult(fetchResult)
	}
}
