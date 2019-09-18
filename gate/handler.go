package gate

import (
	fmt "fmt"
	"strings"

	"github.com/liangdas/mqant/log"
	"github.com/liangdas/mqant/utils"
)

type handler struct {
	gate     Gate
	sessions *utils.BeeMap //use sessionID be key
}

func NewGateHandler(gate Gate) *handler {
	handler := &handler{
		gate:     gate,
		sessions: utils.NewBeeMap(),
	}

	return handler
}

func (h *handler) OnDestory() {
	for _, v := range h.sessions.Items() {
		v.(Agent).Close()
	}
	h.sessions.DeleteAll()
}

func (h *handler) Update(sessionID string) (reslut Session, err string) {
	agent := h.sessions.Get(sessionID)
	if agent == nil {
		err = "No agent found"
		return
	}
	reslut = agent.(Agent).GetSession()
	return
}

func (h *handler) Bind(Sessionid string, Userid string) (result Session, err string) {
	agent := h.sessions.Get(Sessionid)
	if agent == nil {
		err = "No Sesssion found"
		return
	}
	agent.(Agent).GetSession().SetUserid(Userid)

	if h.gate.GetStorageHandler() != nil && agent.(Agent).GetSession().GetUserid() != "" {
		data, err := h.gate.GetStorageHandler().Query(Userid)
		if err == nil && data != nil {
			//上一次保存的連接
			imSession, err := h.gate.NewSession(data)
			if err == nil {
				if agent.(Agent).GetSession().GetSettings() == nil {
					agent.(Agent).GetSession().SetSettings(imSession.GetSettings())
				} else {
					settings := imSession.GetSettings()
					if settings != nil {
						for k, v := range settings {
							if _, ok := agent.(Agent).GetSession().GetSettings()[k]; ok {
								//不替換
							} else {
								agent.(Agent).GetSession().GetSettings()[k] = v
							}
						}
					}
					//數據持久化
					h.gate.GetStorageHandler().Storage(Userid, agent.(Agent).GetSession())
				}
			} else {
				//解析持久化數據失敗
			}
		}
	}
	result = agent.(Agent).GetSession()
	return
}

func (h *handler) IsConnect(Sessionid string, Userid string) (bool, string) {
	for _, agent := range h.sessions.Items() {
		if agent.(Agent).GetSession().GetUserid() == Userid {
			return !agent.(Agent).IsClosed(), ""
		}
	}
	return false, fmt.Sprintf("The gateway did not find the corresponding userId 【%s】", Userid)
}

func (h *handler) UnBind(sessionID string) (result Session, err string) {
	agent := h.sessions.Get(sessionID)
	if agent == nil {
		err = "No Session found"
		return
	}
	agent.(Agent).GetSession().SetUserid("")
	result = agent.(Agent).GetSession()
	return
}

func (h *handler) Push(sessionID string, settings map[string]string) (result Session, err string) {
	agent := h.sessions.Get(sessionID)
	if agent == nil {
		err = "No Session found"
		return
	}
	agent.(Agent).GetSession().SetSettings(settings)
	result = agent.(Agent).GetSession()
	if h.gate.GetStorageHandler() != nil && agent.(Agent).GetSession().GetUserid() != "" {
		err := h.gate.GetStorageHandler().Storage(agent.(Agent).GetSession().GetUserid(), agent.(Agent).GetSession())
		if err != nil {
			log.Warning("gate session storage failure : %s", err.Error())
		}
	}

	return
}

func (h *handler) Set(sessionID string, key string, value string) (result Session, err string) {
	agent := h.sessions.Get(sessionID)
	if agent == nil {
		err = "No Session found"
		return
	}
	agent.(Agent).GetSession().GetSettings()[key] = value
	result = agent.(Agent).GetSession()

	if h.gate.GetStorageHandler() != nil && agent.(Agent).GetSession().GetUserid() != "" {
		err := h.gate.GetStorageHandler().Storage(agent.(Agent).GetSession().GetUserid(), agent.(Agent).GetSession())
		if err != nil {
			log.Error("gate session storage failure : %s", err.Error())
		}
	}
	return
}

func (h *handler) Remove(sessionID string, key string) (reslut interface{}, err string) {
	agent := h.sessions.Get(sessionID)
	if agent == nil {
		err = "No Session found"
		return
	}
	delete(agent.(Agent).GetSession().GetSettings(), key)
	reslut = agent.(Agent).GetSession()

	if h.gate.GetStorageHandler() != nil && agent.(Agent).GetSession().GetUserid() != "" {
		err := h.gate.GetStorageHandler().Storage(agent.(Agent).GetSession().GetUserid(), agent.(Agent).GetSession())
		if err != nil {
			log.Error("gate session storage failure :%s", err.Error())
		}
	}

	return
}

func (h *handler) Send(sessionID string, data interface{}) (err string) {
	agent := h.sessions.Get(sessionID)
	if agent == nil {
		err = "No session found"
		return
	}
	agent.(Agent).WriteMsg(data)
	return
}

func (h *handler) SendBatch(sessionIDStr string, data interface{}) (int64, string) {
	sessionIDs := strings.Split(sessionIDStr, ",")
	var count int64 = 0
	for _, sessionID := range sessionIDs {
		agent := h.sessions.Get(sessionID)
		if agent == nil {
			continue
		}
		agent.(Agent).WriteMsg(data)
		count++
	}
	return count, ""
}

func (h *handler) BroadCast(data interface{}) int64 {
	var count int64 = 0
	for _, agent := range h.sessions.Items() {
		agent.(Agent).WriteMsg(data)
		count++
	}
	return count
}

func (h *handler) Close(sessionID string) error {
	agent := h.sessions.Get(sessionID)
	if agent == nil {
		return fmt.Errorf("No session found")
	}
	agent.(Agent).Close()
	return nil
}
