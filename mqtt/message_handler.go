package mqtt

import (
    "github.com/eclipse/paho.golang/paho"
    "go-linux-iot/logx"
    "go-linux-iot/shared"
)

type Handler struct {
}

func (handler *Handler) Handle(msg *paho.Publish) {
    m := make(map[string]any)
    m["msg"] = msg
    m["payload"] = string(msg.Payload)
    shared.LOGGER.Info("收到平台下行消息", logx.MsgInfo(m))

}
