package gate

import "net"

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
	IsClosed() bool
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
	Update() (err error)
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
	Close() (err error)
}

/**
Session信息持久化
*/
type StorageHandler interface {
	/**
	存储用户的Session信息
	Session Bind Userid以后每次设置 settings都会调用一次Storage
	*/
	Storage(Userid string, session Session) (err error)
	/**
	强制删除Session信息
	*/
	Delete(Userid string) (err error)
	/**
	获取用户Session信息
	Bind Userid时会调用Query获取最新信息
	*/
	Query(Userid string) (data []byte, err error)
	/**
	用户心跳,一般用户在线时1s发送一次
	可以用来延长Session信息过期时间
	*/
	Heartbeat(Userid string)
}

type SessionLearner interface {
	Connect(a Session)    //当连接建立  并且MQTT协议握手成功
	DisConnect(a Session) //当连接关闭	或者客户端主动发送MQTT DisConnect命令
}

type AgentLearner interface {
	Connect(a Agent)    //当连接建立  并且MQTT协议握手成功
	DisConnect(a Agent) //当连接关闭	或者客户端主动发送MQTT DisConnect命令
}
