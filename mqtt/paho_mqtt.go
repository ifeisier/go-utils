package mqtt

import (
    "fmt"
    "go-linux-iot/logx"
    "go-linux-iot/shared"
)

type PahoMQTTLog struct {
    Prefix     string
    IsErrorLog bool
}

func (l PahoMQTTLog) Println(v ...interface{}) {
    m := make(map[string]any)
    m["log"] = v
    if l.IsErrorLog {
        shared.LOGGER.Error(l.Prefix, logx.MsgInfo(m))
    } else {
        shared.LOGGER.Debug(l.Prefix, logx.MsgInfo(m))
    }
}

func (l PahoMQTTLog) Printf(format string, v ...interface{}) {
    m := make(map[string]any)
    m["log"] = fmt.Sprintf(format, v)
    if l.IsErrorLog {
        shared.LOGGER.Error(l.Prefix, logx.MsgInfo(m))
    } else {
        shared.LOGGER.Debug(l.Prefix, logx.MsgInfo(m))
    }
}
