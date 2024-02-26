package mqtt

import (
    "errors"
    "github.com/eclipse/paho.golang/autopaho"
    "go-linux-iot/logx"
    "go-linux-iot/shared"
)

type Handler struct {
}

func (handler *Handler) Handle(received autopaho.PublishReceived) (bool, error) {
    m := make(map[string]any)
    m["msg"] = received.Packet
    shared.LOGGER.Info("收到平台下行消息", logx.MsgInfo(m))

    // 多个 Handle 是按照顺序执行的，返回值会传递给下一个 Handle。
    //
    // 返回 true 表示消息已经被处理了，下一个 Handle 可以通过 AlreadyHandled 属性判断消息是否已经被处理，
    // 这样做可以避免处理相同的消息。
    //
    // 返回 error 不会立即报错，而是将错误消息传递给下一个 Handle，可以通过 Errs 切片获取。
    return true, errors.New("ads")
}
