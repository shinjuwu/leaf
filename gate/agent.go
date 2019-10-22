package gate

import (
	"net"
	"reflect"

	"github.com/shinjuwu/leaf/util"

	"github.com/shinjuwu/leaf/log"
	"github.com/shinjuwu/leaf/network"
)

type agent struct {
	agentID  int64
	tableID  int
	conn     network.Conn
	gate     *Gate
	userData interface{}
	session  Session
	isclose  bool
}

func (a *agent) Run() {
	for {
		data, err := a.conn.ReadMsg()
		if err != nil {
			log.Debug("read message: %v", err)
			break
		}
		if a.GetSession() == nil {
			a.session, err = NewSessionByMap(a.gate.AgentChanRPC, map[string]interface{}{
				"Sessionid": util.GenerateID().String(),
				"Network":   a.conn.RemoteAddr().Network(),
				"IP":        a.conn.RemoteAddr().String(),
				"Settings":  make(map[string]string),
			})
		}
		a.gate.GetAgentLearner().Connect(a)
		if a.gate.Processor != nil {
			msg, err := a.gate.Processor.Unmarshal(data)
			if err != nil {
				log.Debug("unmarshal message error: %v", err)
				break
			}
			err = a.gate.Processor.Route(msg, a)
			if err != nil {
				log.Debug("route message error: %v", err)
				break
			}
		}
	}
}

func (a *agent) OnClose() {
	a.isclose = true
	a.gate.GetAgentLearner().DisConnect(a) //发送连接断开的事件
}

func (a *agent) WriteMsg(msg interface{}) {
	if a.gate.Processor != nil {
		data, err := a.gate.Processor.Marshal(msg)
		if err != nil {
			log.Error("marshal message %v error: %v", reflect.TypeOf(msg), err)
			return
		}
		err = a.conn.WriteMsg(data...)
		if err != nil {
			log.Error("write message %v error: %v", reflect.TypeOf(msg), err)
		}
	}
}

func (a *agent) WriteMsgByte(data [][]byte) {
	err := a.conn.WriteMsg(data...)
	if err != nil {
		log.Error("write message %v error: %v", data, err)
	}
}

func (a *agent) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *agent) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

func (a *agent) Close() {
	a.conn.Close()
}

func (a *agent) Destroy() {
	a.conn.Destroy()
}

func (a *agent) UserData() interface{} {
	return a.userData
}

func (a *agent) SetUserData(data interface{}) {
	a.userData = data
}

func (a *agent) SetAgentID(id int64) {
	a.agentID = id
}

func (a *agent) GetAgentID() int64 {
	return a.agentID
}

func (a agent) IsClosed() bool {
	return a.isclose
}
func (a *agent) GetSession() Session {
	return a.session
}
func (a *agent) SetTableID(id int) {
	a.tableID = id
}

func (a *agent) GetTableID() int {
	return a.tableID
}
