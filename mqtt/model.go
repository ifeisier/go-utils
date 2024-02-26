package mqtt

import (
    "github.com/eclipse/paho.golang/autopaho"
    "github.com/eclipse/paho.golang/paho"
    "sync"
)

type Client struct {
    serverURL         []string
    clientID          string
    username          string
    password          []byte
    subscribe         *paho.Subscribe
    clientConfig      paho.ClientConfig
    connectionManager *autopaho.ConnectionManager
    error             error
    once              sync.Once
}

// NewClient 创建一个新的 mqtt 客户端
//
// serverURL: mqtt 服务器地址
//
// subscribe: 订阅的主题
func NewClient(serverURL []string, subscribe *paho.Subscribe) *Client {
    return &Client{
        serverURL: serverURL,
        subscribe: subscribe,
    }
}

// NewClientStr 创建一个新的 mqtt 客户端
//
// serverURL: mqtt 服务器地址
//
// subscribe: 订阅的主题
func NewClientStr(serverURL []string, subscribe []string) *Client {
    options := make([]paho.SubscribeOptions, 0, len(subscribe))
    for _, v := range subscribe {
        options = append(options, paho.SubscribeOptions{Topic: v, QoS: 2})
    }

    return &Client{
        serverURL: serverURL,
        subscribe: &paho.Subscribe{
            Subscriptions: options,
        },
    }
}
