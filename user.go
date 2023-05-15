package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
	}

	//启动监听当前User channel 的goroutine
	go user.ListenMsg()

	return user
}

// ListenMsg 监听当前User channel 的方法， 一旦有消息就直接发送给对端客户端
func (user *User) ListenMsg() {
	for {
		msg := <-user.C

		user.conn.Write([]byte(msg + "\r\n"))
	}
}
