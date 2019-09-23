package gate

import (
	fmt "fmt"

	"github.com/golang/protobuf/proto"
	"github.com/shinjuwu/leaf/chanrpc"
)

type sessionagent struct {
	AgentChanRPC *chanrpc.Server
	session      *SessionImp
	judgeGuest   func(session Session) bool
}

func NewSession(chanRpc *chanrpc.Server, data []byte) (Session, error) {
	agent := &sessionagent{AgentChanRPC: chanRpc}
	se := &SessionImp{}
	err := proto.Unmarshal(data, se)
	if err != nil {
		return nil, err
	} // 测试结果
	agent.session = se
	return agent, nil
}

func NewSessionByMap(chanRpc *chanrpc.Server, data map[string]interface{}) (Session, error) {
	agent := &sessionagent{
		AgentChanRPC: chanRpc,
		session:      new(SessionImp),
	}
	err := agent.updateMap(data)
	if err != nil {
		return nil, err
	}
	return agent, nil
}

func (this *sessionagent) GetIP() string {
	return this.session.GetIP()
}

func (this *sessionagent) GetNetwork() string {
	return this.session.GetNetwork()
}

func (this *sessionagent) GetUserid() string {
	return this.session.GetUserid()
}

func (this *sessionagent) GetSessionid() string {
	return this.session.GetSessionid()
}

func (this *sessionagent) GetServerid() string {
	return this.session.GetServerid()
}

func (this *sessionagent) GetSettings() map[string]string {
	return this.session.GetSettings()
}

func (this *sessionagent) SetIP(ip string) {
	this.session.IP = ip
}
func (this *sessionagent) SetNetwork(network string) {
	this.session.Network = network
}
func (this *sessionagent) SetUserid(userid string) {
	this.session.Userid = userid
}
func (this *sessionagent) SetSessionid(sessionid string) {
	this.session.Sessionid = sessionid
}
func (this *sessionagent) SetServerid(serverid string) {
	this.session.Serverid = serverid
}
func (this *sessionagent) SetSettings(settings map[string]string) {
	this.session.Settings = settings
}

func (this *sessionagent) updateMap(s map[string]interface{}) error {
	Userid := s["Userid"]
	if Userid != nil {
		this.session.Userid = Userid.(string)
	}
	IP := s["IP"]
	if IP != nil {
		this.session.IP = IP.(string)
	}
	Network := s["Network"]
	if Network != nil {
		this.session.Network = Network.(string)
	}
	Sessionid := s["Sessionid"]
	if Sessionid != nil {
		this.session.Sessionid = Sessionid.(string)
	}
	Serverid := s["Serverid"]
	if Serverid != nil {
		this.session.Serverid = Serverid.(string)
	}
	Settings := s["Settings"]
	if Settings != nil {
		this.session.Settings = Settings.(map[string]string)
	}
	return nil
}

func (this *sessionagent) update(s Session) error {
	Userid := s.GetUserid()
	this.session.Userid = Userid
	IP := s.GetIP()
	this.session.IP = IP
	Network := s.GetNetwork()
	this.session.Network = Network
	Sessionid := s.GetSessionid()
	this.session.Sessionid = Sessionid
	Serverid := s.GetServerid()
	this.session.Serverid = Serverid
	Settings := s.GetSettings()
	this.session.Settings = Settings
	return nil
}

func (this *sessionagent) Serializable() ([]byte, error) {
	data, err := proto.Marshal((this.session))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (this *sessionagent) Update() (err error) {
	if this.AgentChanRPC == nil {
		err = fmt.Errorf("AgentChanRPC is nil")
		return
	}
	result, err := this.AgentChanRPC.Call1("Update", this.session.Sessionid)
	if err == nil {
		if result != nil {
			//绑定成功,重新更新当前Session
			this.update(result.(Session))
		}
	}
	return
}

func (this *sessionagent) Bind(userid string) string {
	if this.AgentChanRPC == nil {
		return "AgentChanRPC is nil"
	}

	result, err := this.AgentChanRPC.Call1("Bind", this.session.Sessionid, userid)
	if err == nil {
		if result != nil {
			//绑定成功,重新更新当前Session
			this.update(result.(Session))
		}
	}
	return ""
}

func (this *sessionagent) UnBind() string {
	if this.AgentChanRPC == nil {
		return "AgentChanRPC is nil"
	}
	result, err := this.AgentChanRPC.Call1("Unbind", this.session.Sessionid)
	if err == nil {
		if result != nil {
			//绑定成功,重新更新当前Session
			this.update(result.(Session))
		}
	}
	return ""
}

func (this *sessionagent) Push() string {
	if this.AgentChanRPC == nil {
		return "AgentChanRPC is nil"
	}

	result, err := this.AgentChanRPC.Call1("Push", this.session.Sessionid, this.session.Settings)
	if err == nil {
		if result != nil {
			//绑定成功,重新更新当前Session
			this.update(result.(Session))
		}
	}
	return ""
}

func (this *sessionagent) Set(key string, value string) string {
	if this.AgentChanRPC == nil {
		return "AgentChanRPC is nil"
	}
	if this.session.Settings == nil {
		this.session.Settings = map[string]string{}
	}
	this.session.Settings[key] = value
	return ""
}
func (this *sessionagent) SetPush(key string, value string) string {
	if this.AgentChanRPC == nil {
		return "AgentChanRPC is nil"
	}
	if this.session.Settings == nil {
		this.session.Settings = map[string]string{}
	}
	this.session.Settings[key] = value
	return this.Push()
}
func (this *sessionagent) Get(key string) (result string) {
	if this.session.Settings == nil {
		return
	}
	result = this.session.Settings[key]
	return
}

func (this *sessionagent) Remove(key string) string {
	if this.AgentChanRPC == nil {
		return "AgentChanRPC is nil"
	}
	if this.session.Settings == nil {
		this.session.Settings = map[string]string{}
	}
	delete(this.session.Settings, key)
	return ""
}
func (this *sessionagent) Send(id string, data interface{}) string {
	if this.AgentChanRPC == nil {
		return "AgentChanRPC is nil"
	}
	_, err := this.AgentChanRPC.Call1("Send", this.session.Sessionid, data)
	if err != nil {
		return err.Error()
	}
	return ""
}

func (this *sessionagent) SendBatch(Sessionids string, data interface{}) string {
	if this.AgentChanRPC == nil {
		return "AgentChanRPC is nil"
	}

	_, err := this.AgentChanRPC.Call1("SendBatch", Sessionids, data)
	if err != nil {
		return err.Error()
	}
	return ""
}

func (this *sessionagent) IsConnect(userId string) string {
	if this.AgentChanRPC == nil {
		return "AgentChanRPC is nil"
	}
	_, err := this.AgentChanRPC.Call1("IsConnect", this.session.Sessionid, userId)
	if err != nil {
		return err.Error()
	}
	return ""
}

func (this *sessionagent) SendNR(id string, data interface{}) string {
	if this.AgentChanRPC == nil {
		return "AgentChanRPC is nil"
	}
	err := this.AgentChanRPC.Call0("Send", this.session.Sessionid, data)
	if err != nil {
		return err.Error()
	}
	return ""
}

func (this *sessionagent) Close() error {
	if this.AgentChanRPC == nil {
		return fmt.Errorf("AgentChanRPC is nil")
	}
	_, err := this.AgentChanRPC.Call1("Close", this.session.Sessionid)
	return err
}

//是否是访客(未登录) ,默认判断规则为 userId==""代表访客
func (this *sessionagent) IsGuest() bool {
	if this.judgeGuest != nil {
		return this.judgeGuest(this)
	}
	if this.GetUserid() == "" {
		return true
	} else {
		return false
	}
}

//设置自动的访客判断函数,记得一定要在全局的时候设置这个值,以免部分模块因为未设置这个判断函数造成错误的判断
func (this *sessionagent) JudgeGuest(judgeGuest func(session Session) bool) {
	this.judgeGuest = judgeGuest
}
