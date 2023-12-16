package mqtt_util

import (
	"context"
	"errors"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"math/rand"
	"reflect"
	"strconv"
	"time"
)

// DefaultConnectionManager 获取空的 MQTT 连接
func DefaultConnectionManager() **autopaho.ConnectionManager {
	ptr := &autopaho.ConnectionManager{}
	return &ptr
}

// Create 创建 mqtt 的 autopaho.ConnectionManager,
//
// config 是 mqtt 的配置
//
// handler 是 mqtt 收到消息的回调
//
// connectionManager 是 mqtt 的连接实例,
// 通过这个二级指针返回新的 mqtt 连接实例, 还会关闭旧的连接实例。
func Create(config *ClientConfig, handler func(*paho.Publish), connectionManager **autopaho.ConnectionManager) (err error) {
	if config.ClientID == "" {
		config.ClientID = "api-"
		int31n := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(16)
		for _, v := range int31n {
			config.ClientID += strconv.Itoa(v)
		}
	}

	newManager, err := NewClientBuild().AddServerURL(config.ServerURL...).
		AddSubscribe(config.Subscribe...).SetClientID(config.ClientID).
		SetUsernamePassword(config.UserName, config.PassWord).
		SetMessageCallback(handler).
		BuildAndConnection()
	if err != nil {
		return
	}

	if connectionManager == nil {
		return errors.New("connectionManager 不能是 nil")
	}

	if !reflect.DeepEqual(*connectionManager, &autopaho.ConnectionManager{}) {
		defer func(oldConn *autopaho.ConnectionManager) {
			err = oldConn.Disconnect(context.Background())
		}(*connectionManager)
	}
	*connectionManager = newManager

	return
}
