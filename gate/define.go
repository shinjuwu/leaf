package gate

import (
	"net"
)

type Agent interface {
	WriteMsg(msg interface{})
	WriteMsgByte(data [][]byte)
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
	SetTableID(id int)
	GetTableID() int
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
	Bind(Userid string) string
	UnBind() string
	Push() string
	Set(key string, value string) string
	SetPush(key string, value string) string //设置值以后立即推送到gate网关
	Get(key string) string
	Remove(key string) string
	Send(id string, data interface{}) string
	SendNR(id string, data interface{}) string
	SendBatch(Sessionids string, data interface{}) string //想该客户端的网关批量发送消息
	//查询某一个userId是否连接中，这里只是查询这一个网关里面是否有userId客户端连接，如果有多个网关就需要遍历了
	IsConnect(Userid string) string
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

/**
net代理服务 处理器
*/
type GateHandler interface {
	Bind(args []interface{}) interface{}       //Bind the session with the the Userid.
	UnBind(args []interface{}) interface{}     //UnBind the session with the the Userid.
	Set(args []interface{}) interface{}        //Set values (one or many) for the session.
	Remove(args []interface{}) interface{}     //Remove value from the session.
	Push(args []interface{}) (Session, string) //推送信息给Session
	Send(args []interface{}) interface{}       //Send message
	SendBatch(args []interface{}) interface{}  //批量发送
	BroadCast(args []interface{}) interface{}  //广播消息给网关所有在连客户端
	//查询某一个userId是否连接中，这里只是查询这一个网关里面是否有userId客户端连接，如果有多个网关就需要遍历了
	IsConnect(args []interface{}) interface{}
	Close(args []interface{}) interface{}  //主动关闭连接
	Update(args []interface{}) interface{} //更新整个Session 通常是其他模块拉取最新数据
	OnDestory()                            //退出事件,主动关闭所有的连接
}
