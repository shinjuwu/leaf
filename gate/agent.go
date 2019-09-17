package gate

import (
	"net"
)

type Agent interface {
	WriteMsg(msg interface{})
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	Destroy()
	UserData() interface{}
	SetUserData(data interface{})
	GetSession() Session
	SetAgentID(id int64)
	GetAgentID() int64
}

type Session interface {
	GetIP() string
	GetNetwork() string
	GetUserid() string
	GetSessionid() string
	GetServerid() string
	GetSettings() map[string]string
	SetIP(ip string)
	SetNetwork(network string)
	SetUserid(userid string)
	SetSessionid(sessionid string)
	SetServerid(serverid string)
	SetSettings(settings map[string]string)
	Serializable() ([]byte, error)
	Update() (err string)
	Bind(Userid string) (err error)
	UnBind() (err error)
	Push() (err error)
	Set(key string, value string) (err error)
	SetPush(key string, value string) (err error) //设置值以后立即推送到gate网关
	Get(key string) (result string)
	Remove(key string) (err error)
	Send(topic string, body []byte) (err error)
	SendNR(topic string, body []byte) (err error)
	SendBatch(Sessionids string, topic string, body []byte) (int64, error) //想该客户端的网关批量发送消息
	//查询某一个userId是否连接中，这里只是查询这一个网关里面是否有userId客户端连接，如果有多个网关就需要遍历了
	IsConnect(Userid string) (result bool, err error)
	//是否是访客(未登录) ,默认判断规则为 userId==""代表访客
	IsGuest() bool
	//设置自动的访客判断函数,记得一定要在全局的时候设置这个值,以免部分模块因为未设置这个判断函数造成错误的判断
	JudgeGuest(judgeGuest func(session Session) bool)
	Close() (err string)
}
