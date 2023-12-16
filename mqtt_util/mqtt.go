package mqtt_util

import (
	"context"
	"fmt"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"log"
	"net/url"
)

type ConnectionSuccess func(*autopaho.ConnectionManager, *paho.Connack)
type ConnectionError func(err error)
type ClientError func(err error)
type ServerDisconnect func(*paho.Disconnect)

type ClientConfig struct {
	ServerURL []string `yaml:"serverURL" json:"serverURL"`
	Subscribe []string `yaml:"subscribe" json:"subscribe"`
	ClientID  string   `yaml:"clientID" json:"clientID"`
	UserName  string   `yaml:"userName" json:"userName"`
	PassWord  string   `yaml:"passWord" json:"passWord"`
}

type ClientBuild struct {
	ClientConfig
	MessageCallback   paho.MessageHandler
	ConnectionSuccess ConnectionSuccess
	ConnectionError   ConnectionError
	ClientError       ClientError
	ServerDisconnect  ServerDisconnect
}

func NewClientBuild() *ClientBuild {
	return &ClientBuild{
		ConnectionSuccess: func(manager *autopaho.ConnectionManager, connack *paho.Connack) {
			log.Println("成功连接到 MQTT 服务端")
		},
		ConnectionError: func(err error) {
			log.Println("连接 MQTT 服务端失败")
		},
		ClientError: func(err error) {
			log.Println("其它原因导致客户端连接断开")
		},
		ServerDisconnect: func(disconnect *paho.Disconnect) {
			log.Println("收到服务端的断开连接命令")
		},
		MessageCallback: func(publish *paho.Publish) {
			log.Println("收到 MQTT 消息")
		},
	}
}

func (build *ClientBuild) BuildAndConnection() (connectionManager *autopaho.ConnectionManager, err error) {
	if len(build.ServerURL) == 0 {
		err = fmt.Errorf("没有指定 MTTT 服务端")
	}

	urls := make([]*url.URL, 0, len(build.ServerURL))
	for _, v := range build.ServerURL {
		su, _ := url.Parse(v)
		urls = append(urls, su)
	}

	cliCfg := autopaho.ClientConfig{
		BrokerUrls: urls,
		KeepAlive:  30,
		OnConnectionUp: func(manager *autopaho.ConnectionManager, connack *paho.Connack) {
			subscribeLen := len(build.Subscribe)
			if subscribeLen != 0 {
				sos := make([]paho.SubscribeOptions, 0, subscribeLen)
				for _, v := range build.Subscribe {
					sos = append(sos, paho.SubscribeOptions{Topic: v, QoS: 2})
				}

				_, err = connectionManager.Subscribe(context.Background(), &paho.Subscribe{Subscriptions: sos})
				if err != nil {
					return
				}
			}

			build.ConnectionSuccess(manager, connack)
		},
		OnConnectError: build.ConnectionError,
		Debug:          paho.NOOPLogger{},
		ClientConfig: paho.ClientConfig{
			ClientID:           build.ClientID,
			OnClientError:      build.ClientError,
			OnServerDisconnect: build.ServerDisconnect,
		},
	}

	if build.UserName == "" || build.PassWord == "" {
		err = fmt.Errorf("没有指定 MQTT 用户名和密码")
	}
	cliCfg.SetUsernamePassword(build.UserName, []byte(build.PassWord))
	cliCfg.ClientConfig.Router = paho.NewSingleHandlerRouter(build.MessageCallback)

	connectionManager, err = autopaho.NewConnection(context.Background(), cliCfg)
	if err != nil {
		return
	}

	go func(cm *autopaho.ConnectionManager, cb *ClientBuild) {
		err = cm.AwaitConnection(context.Background())
		if err != nil {
			cb.ConnectionError(err)
			return
		}
	}(connectionManager, build)

	return
}

// AddServerURL 添加服务器
func (build *ClientBuild) AddServerURL(url ...string) *ClientBuild {
	build.ServerURL = append(build.ServerURL, url...)
	return build
}

// AddSubscribe 添加订阅
func (build *ClientBuild) AddSubscribe(subscribe ...string) *ClientBuild {
	build.Subscribe = append(build.Subscribe, subscribe...)
	return build
}

// SetClientID 设置客户端ID
func (build *ClientBuild) SetClientID(clientID string) *ClientBuild {
	build.ClientID = clientID
	return build
}

// SetUsernamePassword 设置用户名和密码
func (build *ClientBuild) SetUsernamePassword(un, pwd string) *ClientBuild {
	build.UserName = un
	build.PassWord = pwd
	return build
}

// SetMessageCallback 设置消息回调
func (build *ClientBuild) SetMessageCallback(m paho.MessageHandler) *ClientBuild {
	build.MessageCallback = m
	return build
}

// SetConnectionSuccess 成功连接到 MQTT 服务端
func (build *ClientBuild) SetConnectionSuccess(connectionSuccess ConnectionSuccess) *ClientBuild {
	build.ConnectionSuccess = connectionSuccess
	return build
}

// SetConnectionError 连接 MQTT 服务端失败
func (build *ClientBuild) SetConnectionError(connectionError ConnectionError) *ClientBuild {
	build.ConnectionError = connectionError
	return build
}

// SetClientError 其它原因导致客户端连接断开
func (build *ClientBuild) SetClientError(clientError ClientError) *ClientBuild {
	build.ClientError = clientError
	return build
}

// SetServerDisconnect 收到服务端的断开连接命令
func (build *ClientBuild) SetServerDisconnect(serverDisconnect ServerDisconnect) *ClientBuild {
	build.ServerDisconnect = serverDisconnect
	return build
}
