package gate

import "github.com/liangdas/mqant/utils"

type handler struct {
	gate     Gate
	sessions *utils.BeeMap
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

func (h *handler) Bind(userID string) error {

}

func (h *handler) IsConnect(agentID int64) (bool, string) {
	for _, agent := range h.sessions.Items() {

	}
}
