package gate

import (
	"time"

	"github.com/shinjuwu/leaf/chanrpc"
	"github.com/shinjuwu/leaf/network"
)

type Gate struct {
	MaxConnNum      int
	PendingWriteNum int
	MaxMsgLen       uint32
	Processor       network.Processor
	CustomProcessor network.Processor
	AgentChanRPC    *chanrpc.Server

	// websocket
	WSAddr      string
	HTTPTimeout time.Duration
	CertFile    string
	KeyFile     string

	// tcp
	TCPAddr      string
	LenMsgLen    int
	LittleEndian bool

	//extension
	handler        GateHandler
	agentLearner   AgentLearner
	sessionLearner SessionLearner
	storage        StorageHandler
	//tracing        TracingHandler
}

func (gate *Gate) OnInit() {
	handler := NewGateHandler(*gate)

	gate.agentLearner = handler
	gate.handler = handler
	gate.AgentChanRPC.Register("Update", gate.handler.Update)
	gate.AgentChanRPC.Register("Bind", gate.handler.Bind)
	gate.AgentChanRPC.Register("UnBind", gate.handler.UnBind)
	gate.AgentChanRPC.Register("Push", gate.handler.Push)
	gate.AgentChanRPC.Register("Set", gate.handler.Set)
	gate.AgentChanRPC.Register("Remove", gate.handler.Remove)
	gate.AgentChanRPC.Register("Send", gate.handler.Send)
	gate.AgentChanRPC.Register("SendBatch", gate.handler.SendBatch)
	gate.AgentChanRPC.Register("BroadCast", gate.handler.BroadCast)
	gate.AgentChanRPC.Register("IsConnect", gate.handler.IsConnect)
	gate.AgentChanRPC.Register("Close", gate.handler.Close)
}

func (gate *Gate) Run(closeSig chan bool) {
	var wsServer *network.WSServer
	if gate.WSAddr != "" {
		wsServer = new(network.WSServer)
		wsServer.Addr = gate.WSAddr
		wsServer.MaxConnNum = gate.MaxConnNum
		wsServer.PendingWriteNum = gate.PendingWriteNum
		wsServer.MaxMsgLen = gate.MaxMsgLen
		wsServer.HTTPTimeout = gate.HTTPTimeout
		wsServer.CertFile = gate.CertFile
		wsServer.KeyFile = gate.KeyFile
		wsServer.NewAgent = func(conn *network.WSConn) network.Agent {
			a := &agent{conn: conn, gate: gate}
			if gate.AgentChanRPC != nil {
				gate.AgentChanRPC.Go("NewAgent", a)
			}
			return a
		}
	}

	var tcpServer *network.TCPServer
	if gate.TCPAddr != "" {
		tcpServer = new(network.TCPServer)
		tcpServer.Addr = gate.TCPAddr
		tcpServer.MaxConnNum = gate.MaxConnNum
		tcpServer.PendingWriteNum = gate.PendingWriteNum
		tcpServer.LenMsgLen = gate.LenMsgLen
		tcpServer.MaxMsgLen = gate.MaxMsgLen
		tcpServer.LittleEndian = gate.LittleEndian
		tcpServer.NewAgent = func(conn *network.TCPConn) network.Agent {
			a := &agent{conn: conn, gate: gate}
			if gate.AgentChanRPC != nil {
				gate.AgentChanRPC.Go("NewAgent", a)
			}
			return a
		}
	}

	if wsServer != nil {
		wsServer.Start()
	}
	if tcpServer != nil {
		tcpServer.Start()
	}
	<-closeSig
	if wsServer != nil {
		wsServer.Close()
	}
	if tcpServer != nil {
		tcpServer.Close()
	}
}

func (gate *Gate) OnDestroy() {
}

func (this *Gate) GetStorageHandler() (storage StorageHandler) {
	return this.storage
}

/**
设置Session信息持久化接口
*/
func (this *Gate) SetStorageHandler(storage StorageHandler) error {
	this.storage = storage
	return nil
}

func (this *Gate) NewSession(data []byte) (Session, error) {
	return NewSession(this.AgentChanRPC, data)
}

/**
设置客户端连接和断开的监听器
*/
func (this *Gate) SetSessionLearner(sessionLearner SessionLearner) error {
	this.sessionLearner = sessionLearner
	return nil
}

func (this *Gate) GetSessionLearner() SessionLearner {
	return this.sessionLearner
}

func (this *Gate) GetAgentLearner() AgentLearner {
	return this.agentLearner
}
