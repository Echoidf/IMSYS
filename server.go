package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int
	// 在线用户的列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	//消息广播的channel
	Message chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

// ListenMessage 监听server.Message广播消息channel的goroutine,一旦有消息就发送给全部的在线user
func (server *Server) ListenMessage() {
	for {
		msg := <-server.Message
		if msg == "" {
			continue
		}
		server.mapLock.Lock()
		for _, cli := range server.OnlineMap {
			cli.C <- msg
		}
		server.mapLock.Unlock()
	}
}

// BroadCast 广播消息
func (server *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	server.Message <- sendMsg
}

func (server *Server) Handler(conn net.Conn) {
	//...当前连接的业务
	fmt.Println("连接建立成功")

	user := NewUser(conn)
	//用户上线

	//加入onlineMap
	server.mapLock.Lock()
	server.OnlineMap[user.Name] = user
	server.mapLock.Unlock()

	//进行广播--即将消息写入server的channel
	server.BroadCast(user, "已上线")

	//当前handler阻塞
	select {}
}

func (server *Server) Start() {
	//socket listener
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		fmt.Println("net listener err:", err)
	}

	//close listener socket
	defer listener.Close()

	go server.ListenMessage()

	//accept
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}
		//do handler
		go server.Handler(conn)
	}

}
