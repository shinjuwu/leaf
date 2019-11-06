package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/shinjuwu/leaf/chanrpc"

	"github.com/shinjuwu/leaf/log"
)

type NeooneProcessor struct {
	msgInfo map[string]*MsgInfo
}

func NewNeooneProcessor() *NeooneProcessor {
	p := new(NeooneProcessor)
	p.msgInfo = make(map[string]*MsgInfo)
	return p
}

//Register is a extension for old Neoone msg struct. In common, don't use this func to register your msg in new project.
func (p *NeooneProcessor) Register(msgID string, msg interface{}) {
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		log.Fatal("json message pointer required")
	}

	if _, ok := p.msgInfo[msgID]; ok {
		log.Fatal("Message %v is already registered", msgID)
	}

	i := new(MsgInfo)
	i.msgType = msgType
	p.msgInfo[msgID] = i
}

//RegisterJSON is used to pure json message
func (p *NeooneProcessor) RegisterJSON(msg interface{}) string {
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		log.Fatal("json message pointer required")
	}
	msgID := msgType.Elem().Name()
	if msgID == "" {
		log.Fatal("unnamed json message")
	}
	if _, ok := p.msgInfo[msgID]; ok {
		log.Fatal("message %v is already registered", msgID)
	}

	i := new(MsgInfo)
	i.msgType = msgType
	p.msgInfo[msgID] = i
	return msgID
}

//SetRouter is a extension for old Neoone msg struct. In common, don't use this func to register your msg in new project.
func (p *NeooneProcessor) SetRouter(msgID string, msgRouter *chanrpc.Server) {
	i, ok := p.msgInfo[msgID]
	if !ok {
		log.Fatal("message %v not registered", msgID)
	}
	i.msgRouter = msgRouter
}

//SetRouterJSON is Used for pure json message
func (p *NeooneProcessor) SetRouterJSON(msg interface{}, msgRouter *chanrpc.Server) {
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		log.Fatal("json message pointer required")
	}
	msgID := msgType.Elem().Name()
	i, ok := p.msgInfo[msgID]
	if !ok {
		log.Fatal("message %v not registered", msgID)
	}

	i.msgRouter = msgRouter
}

//SetHandler is a extension for old Neoone msg struct. In common, don't use this func to register your msg in new project.
func (p *NeooneProcessor) SetHandler(msgID string, msgHandler MsgHandler) {
	i, ok := p.msgInfo[msgID]
	if !ok {
		log.Fatal("message %v not registered", msgID)
	}
	i.msgHandler = msgHandler
}

//SetRawHandler is a extension for old Neoone msg struct. In common, don't use this func to register your msg in new project.
func (p *NeooneProcessor) SetRawHandler(msgID string, msgRawHandler MsgHandler) {
	i, ok := p.msgInfo[msgID]
	if !ok {
		log.Fatal("message %v not registered", msgID)
	}

	i.msgRawHandler = msgRawHandler
}

func (p *NeooneProcessor) Route(msgwithID interface{}, userData interface{}) error {
	// raw
	if msgRaw, ok := msgwithID.(MsgRaw); ok {
		i, ok := p.msgInfo[msgRaw.msgID]
		if !ok {
			return fmt.Errorf("message %v not registered", msgRaw.msgID)
		}
		if i.msgRawHandler != nil {
			i.msgRawHandler([]interface{}{msgRaw.msgID, msgRaw.msgRawData, userData})
		}
		return nil
	}
	msgType := reflect.TypeOf(msgwithID)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		return errors.New("json message pointer required")
	}
	msgID := msgType.Elem().Name()
	fmt.Println(msgID)
	msgMap := msgwithID.(map[string]interface{})
	if len(msgMap) != 1 {
		return fmt.Errorf("invaild msg %v", msgMap)
	}
	for msgID, msg := range msgMap {
		// json
		msgType := reflect.TypeOf(msg)
		if msgType == nil || msgType.Kind() != reflect.Ptr {
			return errors.New("json message pointer required")
		}
		i, ok := p.msgInfo[msgID]
		if !ok {
			return fmt.Errorf("message %v not registered", msgID)
		}
		if i.msgHandler != nil {
			i.msgHandler([]interface{}{msg, userData})
		}
		if i.msgRouter != nil {
			i.msgRouter.Go(msgType, msg, userData)
		}
		return nil
	}
	return fmt.Errorf("not register message %v", msgMap)
}

func (p *NeooneProcessor) Unmarshal(data []byte) (interface{}, error) {
	var m map[string]interface{}
	err := json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}
	if id, ok := m["cmd"]; ok {
		msgID := id.(string)
		i, ok := p.msgInfo[msgID]
		if !ok {
			return nil, fmt.Errorf("message %v not register", msgID)
		}
		if i.msgRawHandler != nil {
			return MsgRaw{msgID, data}, nil
		} else {
			msg := reflect.New(i.msgType.Elem()).Interface()
			msgwithID := map[string]interface{}{
				msgID: msg,
			}
			if cmdData, ok := m["data"]; ok {
				return msgwithID, json.Unmarshal([]byte(cmdData.(string)), &msg)
			} else {
				return nil, errors.New("invalid json data")
			}
		}
	} else {
		if len(m) != 1 {
			return nil, errors.New("invalid json data")
		}

		for msgID, v := range m {
			i, ok := p.msgInfo[msgID]
			if !ok {
				return nil, fmt.Errorf("message %v not registered", msgID)
			}

			// msg
			if i.msgRawHandler != nil {
				return MsgRaw{msgID, data}, nil
			} else {
				msg := reflect.New(i.msgType.Elem()).Interface()
				return msg, json.Unmarshal([]byte(v.(string)), &msg)
			}
		}
	}
	return nil, errors.New("invalid json data")
}

func (p *NeooneProcessor) Marshal(msg interface{}) ([][]byte, error) {
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		return nil, errors.New("json message pointer required")
	}
	msgID := msgType.Elem().Name()
	if _, ok := p.msgInfo[msgID]; !ok {
		return nil, fmt.Errorf("message %v not registered", msgID)
	}

	// data
	data, err := json.Marshal(msg)
	return [][]byte{data}, err
}
