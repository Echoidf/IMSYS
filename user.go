package main

import "net"

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
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

// Online 用户上线
func (user *User) Online() {
	//加入onlineMap
	user.server.mapLock.Lock()
	user.server.OnlineMap[user.Name] = user
	user.server.mapLock.Unlock()

	//进行广播--即将消息写入server的channel
	user.server.BroadCast(user, "已上线")
}

// OffLine 用户下线
func (user *User) OffLine() {
	//从onlineMap删除
	user.server.mapLock.Lock()
	delete(user.server.OnlineMap, user.Name)
	user.server.mapLock.Unlock()

	//进行广播--即将消息写入server的channel
	user.server.BroadCast(user, "已下线")
}

func (user *User) DoMessage(msg string) {
	user.server.BroadCast(user, msg)
}
